package controllers

import (
	"amplifier/app/awsservices"
	"amplifier/app/db"
	"amplifier/app/jobs"
	"amplifier/app/jobs/job_handlers"
	"amplifier/app/providers"
	"amplifier/app/work"
	"context"
	"time"

	"github.com/revel/revel"
)

var (
	redisManager         *db.AppRedis
	jobEnqueuer          *work.AppJobEnqueuer
	africasTalkingSender *providers.AppAfricasTalkingSender
	sqsConn              *awsservices.SQSClient
)

func init() {
	revel.OnAppStart(InitApp)
}

func InitApp() {
	africasTalkingSender = providers.NewAfricasTalkingSenderWithParameters(
		revel.Config.StringDefault("aft.url", ""),
		revel.Config.StringDefault("aft.user", ""),
		revel.Config.StringDefault("aft.key", ""),
	)

	redisManager = db.NewRedisProvider(&db.RedisConfig{
		IdleTimeout: 2 * time.Minute,
		MaxActive:   1000,
		MaxIdle:     100,
	})

	redisPool := redisManager.RedisPool()

	jobEnqueuer = work.NewJobEnqueuer(redisPool)

	workerPool := work.NewWorkerPool(redisPool, uint(200))

	jobHandlers := setupJobHandlers(
		africasTalkingSender,
		jobEnqueuer,
	)

	workerPool.RegisterJobs(jobHandlers...)

	workerPool.Start(context.Background())

	sqsConn = awsservices.NewSQSClient()

	spinGoRoutines()
}

func setupJobHandlers(
	africasTalkingSender providers.AfricasTalkingSender,
	jobEnqueuer work.JobEnqueuer,
) []jobs.JobHandler {
	workOnATJobHandler := job_handlers.NewATJobHandler(jobEnqueuer)
	workOnATSendJobHandler := job_handlers.NewATSendJobHandler(africasTalkingSender)
	return []jobs.JobHandler{
		workOnATJobHandler,
		workOnATSendJobHandler,
	}
}

func spinGoRoutines() {
	go processATRequests()

	for i := 0; i < 20; i++ {
		go processATSendRequests(i)
	}
}

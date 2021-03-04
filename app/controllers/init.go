package controllers

import (
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
	redisManager *db.AppRedis
	jobEnqueuer  *work.AppJobEnqueuer
)

func init() {
	revel.OnAppStart(InitApp)
}

func InitApp() {
	redisManager = db.NewRedisProvider(&db.RedisConfig{
		IdleTimeout: 2 * time.Minute,
		MaxActive:   1000,
		MaxIdle:     100,
	})

	redisPool := redisManager.RedisPool()

	jobEnqueuer = work.NewJobEnqueuer(redisPool)

	workerPool := work.NewWorkerPool(redisPool, uint(10))

	jobHandlers := setupJobHandlers(jobEnqueuer)

	workerPool.RegisterJobs(jobHandlers...)

	workerPool.Start(context.Background())
}

func setupJobHandlers(jobEnqueuer work.JobEnqueuer) []jobs.JobHandler {
	africasTalkingSender := providers.NewAfricasTalkingSenderWithParameters(
		revel.Config.StringDefault("aft.url", ""),
		revel.Config.StringDefault("aft.user", ""),
		revel.Config.StringDefault("aft.key", ""),
	)

	workOnATJobHandler := job_handlers.NewATJobHandler(jobEnqueuer)
	workOnATSendJobHandler := job_handlers.NewATSendJobHandler(africasTalkingSender)
	return []jobs.JobHandler{
		workOnATJobHandler,
		workOnATSendJobHandler,
	}
}

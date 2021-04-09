package controllers

import (
	"amplifier/app/db"
	"amplifier/app/providers"
	"amplifier/app/tasks"
	"amplifier/app/tasks/job_handlers"
	"amplifier/app/work"
	"context"
	"time"

	"github.com/revel/revel"
)

var (
	redisManager         *db.AppRedis
	jobEnqueuer          *work.AppJobEnqueuer
	africasTalkingSender *providers.AppAfricasTalkingSender
	// sqsConn              *awsservices.SQSClient
)

func init() {
	revel.OnAppStart(InitApp)
	revel.InterceptMethod(App.AddUser, revel.BEFORE)
	revel.InterceptMethod(Credentials.checkUser, revel.BEFORE)
	revel.InterceptMethod(Requests.checkUser, revel.BEFORE)

	revel.TemplateFuncs["formatDate"] = func(theTime time.Time) string {
		timeLocation, err := time.LoadLocation("Africa/Nairobi")
		if err != nil {
			revel.AppLog.Errorf("failed to load Nairobi timezone: %+v", err)
			return theTime.Format("Jan _2 2006 3:04PM")
		}

		return theTime.In(timeLocation).Format("Jan _2 2006 3:04PM")
	}
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

	jobHandlers := setupTaskHandlers(
		africasTalkingSender,
		jobEnqueuer,
	)

	workerPool.RegisterJobs(jobHandlers...)

	workerPool.Start(context.Background())

	// sqsConn = awsservices.NewSQSClient()

	// spinGoRoutines()
}

func setupTaskHandlers(
	africasTalkingSender providers.AfricasTalkingSender,
	jobEnqueuer work.JobEnqueuer,
) []tasks.JobHandler {
	workOnATJobHandler := job_handlers.NewATJobHandler(jobEnqueuer)
	workOnATSendJobHandler := job_handlers.NewATSendJobHandler(africasTalkingSender)
	return []tasks.JobHandler{
		workOnATJobHandler,
		workOnATSendJobHandler,
	}
}

// func spinGoRoutines() {
// 	go processATRequests()

// 	for i := 0; i < 20; i++ {
// 		go processATSendRequests(i)
// 	}
// }

package tasks

import "context"

type JobHandler interface {
	Job() Job
	PerformJob(ctx context.Context, body string) error
}

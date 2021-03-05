package job_handlers

import (
	"amplifier/app/entities"
	"amplifier/app/jobs"
	"amplifier/app/jobs/sms_jobs"
	"amplifier/app/work"
	"context"
	"encoding/json"
	"fmt"

	"github.com/revel/revel"
)

type ATJobHandler struct {
	jobEnqueuer work.JobEnqueuer
}

func NewATJobHandler(jobEnqueuer work.JobEnqueuer) *ATJobHandler {
	return &ATJobHandler{
		jobEnqueuer: jobEnqueuer,
	}
}

func (h *ATJobHandler) Job() jobs.Job {
	return &sms_jobs.ProcessATJob{}
}

func (h *ATJobHandler) PerformJob(
	ctx context.Context,
	body string,
) error {
	var theJob sms_jobs.ProcessATJob
	err := json.Unmarshal([]byte(body), &theJob)
	if err != nil {
		fmt.Printf("error unmarshal task: %v", err)
		return nil
	}

	revel.AppLog.Infof("ATJobHandler job: =[%+v]", theJob)

	if theJob.Multi {
		for _, rec := range theJob.Recipients {
			job := sms_jobs.NewATSendJob(theJob.SenderID, rec["message"], []*entities.SMSRecipient{{
				Phone: rec["phone"],
			}})
			_, err = h.jobEnqueuer.Enqueue(ctx, job)
			if err != nil {
				revel.AppLog.Errorf("failed to enqueue job: %v", err)
				return err
			}
		}
		return nil
	}

	recs := make([]*entities.SMSRecipient, 0)
	for _, rec := range theJob.Recipients {
		recs = append(recs, &entities.SMSRecipient{
			Phone: rec["phone"],
		})
	}

	job := sms_jobs.NewATSendJob(theJob.SenderID, theJob.Message, recs)
	_, err = h.jobEnqueuer.Enqueue(ctx, job)
	if err != nil {
		revel.AppLog.Errorf("failed to enqueue job: %v", err)
		return err
	}
	return nil
}

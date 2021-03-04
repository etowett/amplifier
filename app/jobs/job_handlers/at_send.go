package job_handlers

import (
	"amplifier/app/entities"
	"amplifier/app/jobs"
	"amplifier/app/jobs/sms_jobs"
	"amplifier/app/providers"
	"context"
	"encoding/json"
	"fmt"

	"github.com/revel/revel"
)

type ATSendJobHandler struct {
	africasTalkingSender providers.AfricasTalkingSender
}

func NewATSendJobHandler(africasTalkingSender providers.AfricasTalkingSender) *ATSendJobHandler {
	return &ATSendJobHandler{
		africasTalkingSender: africasTalkingSender,
	}
}

func (h *ATSendJobHandler) Job() jobs.Job {
	return &sms_jobs.ProcessATJob{}
}

func (h *ATSendJobHandler) PerformJob(
	ctx context.Context,
	body string,
) error {
	var theJob sms_jobs.ATSendJob
	err := json.Unmarshal([]byte(body), &theJob)
	if err != nil {
		fmt.Printf("error unmarshal task: %v", err)
		return nil
	}
	revel.AppLog.Infof("ATSendJobHandler theJob: %+v", theJob)

	_, err = h.africasTalkingSender.Send(&entities.SendRequest{
		SenderID: theJob.SenderID,
		Message:  theJob.Message,
		To:       theJob.Recipients,
	})
	if err != nil {
		revel.AppLog.Errorf("send to at, err: %v\n", err)
		return err
	}
	return nil
}

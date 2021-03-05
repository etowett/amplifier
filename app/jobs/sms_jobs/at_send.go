package sms_jobs

import (
	"amplifier/app/entities"
	"amplifier/app/jobs"
	"encoding/json"
)

const atSendJobName = "at_send_requests"

type ATSendJob struct {
	SenderID   string                   `json:"sender_id"`
	Message    string                   `json:"message"`
	Recipients []*entities.SMSRecipient `json:"recipients"`
}

func NewATSendJob(
	senderID string,
	smsBody string,
	allRecs []*entities.SMSRecipient,
) *ATSendJob {
	return &ATSendJob{
		SenderID:   senderID,
		Message:    smsBody,
		Recipients: allRecs,
	}
}

func (h *ATSendJob) JobName() string {
	return atSendJobName
}

func (h *ATSendJob) JobBody() (string, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (h *ATSendJob) JobOptions() []jobs.PerformJobOption {
	return []jobs.PerformJobOption{
		jobs.WithMaxConcurrency(50),
		jobs.WithMaxFails(2),
		jobs.WithLowPriority(),
	}
}

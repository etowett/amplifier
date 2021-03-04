package sms_jobs

import (
	"amplifier/app/jobs"
	"encoding/json"
)

const processAtJobName = "at_requests"

type ProcessATJob struct {
	Message    string              `json:"message"`
	SenderID   string              `json:"sender_id"`
	Multi      bool                `json:"multi"`
	Recipients []map[string]string `json:"recipients"`
}

func NewATJob(
	smsBody string,
	senderID string,
	isMulti bool,
	allRecs []map[string]string,
) *ProcessATJob {
	return &ProcessATJob{
		Message:    smsBody,
		SenderID:   senderID,
		Multi:      isMulti,
		Recipients: allRecs,
	}
}

func (h *ProcessATJob) JobName() string {
	return processAtJobName
}

func (h *ProcessATJob) JobBody() (string, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (h *ProcessATJob) JobOptions() []jobs.PerformJobOption {
	return []jobs.PerformJobOption{
		jobs.WithMaxConcurrency(5),
		jobs.WithMaxFails(2),
		jobs.WithLowPriority(),
	}
}

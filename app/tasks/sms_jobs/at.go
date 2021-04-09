package sms_jobs

import (
	"amplifier/app/tasks"
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

func (h *ProcessATJob) JobOptions() []tasks.PerformJobOption {
	return []tasks.PerformJobOption{
		tasks.WithMaxConcurrency(10),
		tasks.WithMaxFails(2),
		tasks.WithLowPriority(),
	}
}

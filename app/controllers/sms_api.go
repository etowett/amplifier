package controllers

import (
	"amplifier/app/entities"
	"fmt"
	"net/http"
	"time"

	"github.com/revel/revel"

	"amplifier/app/forms"
	"amplifier/app/helpers"
	"amplifier/app/jobs/sms_jobs"
)

type SMSApi struct {
	App
}

func (c SMSApi) Aft() revel.Result {
	var status int
	atForm := forms.ATForm{}
	err := c.Params.BindJSON(&atForm)
	if err != nil {
		c.Log.Infof("AT form bind error: %[+v], will use defaults", err)
	}
	c.Log.Infof("Given at request: %+v", atForm)

	reqCount := 10
	if atForm.Count > 0 {
		reqCount = int(atForm.Count)
	}

	smsBody := fmt.Sprintf(
		"Hello from the world to %v at %v",
		reqCount,
		time.Now().String()[8:19],
	)

	if len(atForm.Message) > 0 {
		smsBody = atForm.Message
	}

	isMulti := false
	if atForm.Multi {
		isMulti = true
	}

	job := sms_jobs.NewATJob(smsBody, "Amplifier", isMulti, helpers.GetRecipients(reqCount, isMulti))
	_, err = jobEnqueuer.Enqueue(c.Request.Context(), job)
	if err != nil {
		c.Log.Errorf("AT job enqueue: %v\n", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Success: false,
			Status:  status,
			Message: fmt.Sprintf("could not enqueue to at: %v", err),
		})
	}

	status = http.StatusCreated
	c.Response.SetStatus(status)
	return c.RenderJSON(entities.Response{
		Success: true,
		Status:  status,
		Data: map[string]interface{}{
			"count":   reqCount,
			"message": smsBody,
			"multi":   isMulti,
		},
	})
}

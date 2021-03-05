package controllers

import (
	"amplifier/app/entities"
	"amplifier/app/jobs/sms_jobs"
	"fmt"
	"net/http"
	"time"

	"github.com/revel/revel"

	"amplifier/app/forms"
	"amplifier/app/helpers"
)

type SMSApi struct {
	App
}

func (c *SMSApi) validateATRequest(atForm *forms.ATForm) map[string]interface{} {

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
	return map[string]interface{}{
		"count":   reqCount,
		"message": smsBody,
		"multi":   isMulti,
		"recs":    helpers.GetRecipients(reqCount, isMulti),
	}
}

func (c *SMSApi) sendToAT(
	senderID string,
	message string,
	recs []*entities.SMSRecipient,
) error {
	_, err := africasTalkingSender.Send(&entities.SendRequest{
		SenderID: senderID,
		Message:  message,
		To:       recs,
	})
	return err
}

func (c *SMSApi) Aft() revel.Result {
	var status int
	atForm := &forms.ATForm{}
	err := c.Params.BindJSON(atForm)
	if err != nil {
		c.Log.Infof("AT form bind error: %[+v], will use defaults", err)
	}
	c.Log.Infof("Given AT form request: %+v", atForm)

	data := c.validateATRequest(atForm)

	if data["multi"].(bool) {
		for _, rec := range data["recs"].([]map[string]string) {
			recs := []*entities.SMSRecipient{
				{Phone: rec["phone"]},
			}
			err := c.sendToAT("Amplifier", rec["message"], recs)
			if err != nil {
				c.Log.Errorf("to send to AT: %v", err)
				status = http.StatusInternalServerError
				c.Response.SetStatus(status)
				return c.RenderJSON(entities.Response{
					Success: false,
					Status:  status,
					Message: fmt.Sprintf("failed to send to at: %v", err),
				})
			}
		}
	} else {
		recs := make([]*entities.SMSRecipient, 0)
		for _, rec := range data["recs"].([]map[string]string) {
			recs = append(recs, &entities.SMSRecipient{
				Phone: rec["phone"],
			})
		}

		err := c.sendToAT("Amplifier", data["message"].(string), recs)
		if err != nil {
			c.Log.Errorf("to send to AT: %v", err)
			status = http.StatusInternalServerError
			c.Response.SetStatus(status)
			return c.RenderJSON(entities.Response{
				Success: false,
				Status:  status,
				Message: fmt.Sprintf("failed to send to at: %v", err),
			})
		}
	}

	status = http.StatusCreated
	c.Response.SetStatus(status)
	return c.RenderJSON(entities.Response{
		Success: true,
		Status:  status,
		Message: "message(s) sent successfully.",
		Data: map[string]interface{}{
			"count":   data["count"].(int),
			"message": data["message"].(string),
			"multi":   data["multi"].(bool),
		},
	})
}

func (c *SMSApi) AftRedis() revel.Result {
	var status int
	atForm := &forms.ATForm{}
	err := c.Params.BindJSON(atForm)
	if err != nil {
		c.Log.Infof("AT form bind error: %[+v], will use defaults", err)
	}
	c.Log.Infof("Given at request: %+v", atForm)

	data := c.validateATRequest(atForm)
	c.Log.Infof("at ret data: %+v", data)

	job := sms_jobs.NewATJob(
		data["message"].(string),
		"Amplifier",
		data["multi"].(bool),
		data["recs"].([]map[string]string),
	)
	_, err = jobEnqueuer.Enqueue(c.Request.Context(), job)
	if err != nil {
		c.Log.Errorf("Failed AT job enqueue: %v", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Success: false,
			Status:  status,
			Message: fmt.Sprintf("failed to enqueue request for processing: %v", err),
		})
	}

	status = http.StatusCreated
	c.Response.SetStatus(status)
	return c.RenderJSON(entities.Response{
		Success: true,
		Status:  status,
		Message: "message queued for processing.",
		Data: map[string]interface{}{
			"count":   data["count"].(int),
			"message": data["message"].(string),
			"multi":   data["multi"].(bool),
		},
	})
}

func (c SMSApi) AftSQS() revel.Result {
	var status int
	atForm := &forms.ATForm{}
	err := c.Params.BindJSON(atForm)
	if err != nil {
		c.Log.Infof("AT form bind error: %[+v], will use defaults", err)
	}
	c.Log.Infof("Given at request: %+v", atForm)

	data := c.validateATRequest(atForm)
	c.Log.Infof("sqs not implemented - %+v!", data)

	status = http.StatusCreated
	c.Response.SetStatus(status)
	return c.RenderJSON(entities.Response{
		Success: true,
		Status:  status,
		Message: "message queued for processing.",
		Data: map[string]interface{}{
			"count":   data["count"].(int),
			"message": data["message"].(string),
			"multi":   data["multi"].(bool),
		},
	})
}

package controllers

import (
	"amplifier/app/entities"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/revel/revel"
)

func processATRequests() {
	theQueue, err := sqsConn.GetQueueURL(revel.Config.StringDefault("aft.requests_queue", ""))
	if err != nil {
		revel.AppLog.Errorf("could not get sqs url: %v", err)
		return
	}

	for {
		retResponse, err := sqsConn.GetMessageFromQueue(theQueue.QueueUrl)
		if err != nil {
			revel.AppLog.Errorf("error getmessage from queue: %v", err)
			return
		}

		if len(retResponse.Messages) > 0 {

			processedReceiptHandles := make([]*sqs.DeleteMessageBatchRequestEntry, len(retResponse.Messages))

			for i, mess := range retResponse.Messages {
				revel.AppLog.Infof("message =[%v] == %+v", i, mess.String())

				// TODO: introduce buf channels to make this fast!
				err = doProcessToSend(mess.Body)
				if err != nil {
					revel.AppLog.Errorf("could not doProcessToSend: %v", err)
					return
				}

				processedReceiptHandles[i] = &sqs.DeleteMessageBatchRequestEntry{
					Id:            mess.MessageId,
					ReceiptHandle: mess.ReceiptHandle,
				}
			}

			err = sqsConn.DeleteBatchFromQueue(&sqs.DeleteMessageBatchInput{
				QueueUrl: theQueue.QueueUrl,
				Entries:  processedReceiptHandles,
			})
			if err != nil {
				revel.AppLog.Errorf("failed to delete messages batch from queue: %+v", err)
				return
			}

		}

		if len(retResponse.Messages) == 0 {
			revel.AppLog.Infof("processATRequests:: No sqs messages in at requests, sleeping for 10")
			time.Sleep(time.Second * 10)
		}
	}
}

func doProcessToSend(data *string) error {
	theData := map[string]interface{}{}
	err := json.Unmarshal([]byte(*data), &theData)
	if err != nil {
		return err
	}

	theQueue, err := sqsConn.GetQueueURL(revel.Config.StringDefault("aft.send_queue", ""))
	if err != nil {
		return err
	}

	if theData["multi"].(bool) {
		for _, rec := range theData["recs"].([]interface{}) {
			newData := map[string]interface{}{
				"sender_id": theData["sender_id"].(string),
				"message":   rec.(map[string]interface{})["message"].(string),
				"recs": []*entities.SMSRecipient{{
					Phone: rec.(map[string]interface{})["phone"].(string),
				}},
			}
			err = sqsConn.MessageToQueue(theQueue.QueueUrl, newData)
			if err != nil {
				return err
			}
		}
		return nil
	}

	recs := make([]*entities.SMSRecipient, 0)
	for _, rec := range theData["recs"].([]interface{}) {
		recs = append(recs, &entities.SMSRecipient{
			Phone: rec.(map[string]interface{})["phone"].(string),
		})
	}

	newData := map[string]interface{}{
		"sender_id": theData["sender_id"].(string),
		"message":   theData["message"].(string),
		"recs":      recs,
	}
	err = sqsConn.MessageToQueue(theQueue.QueueUrl, newData)
	if err != nil {
		return err
	}

	return nil
}

func processATSendRequests(num int) {
	revel.AppLog.Infof("starting processATSendRequests num: %v", num)
	theQueue, err := sqsConn.GetQueueURL(revel.Config.StringDefault("aft.send_queue", ""))
	if err != nil {
		revel.AppLog.Errorf("could not get sqs url for aft.send_queue: %v", err)
		return
	}

	for {
		retResponse, err := sqsConn.GetMessageFromQueue(theQueue.QueueUrl)
		if err != nil {
			revel.AppLog.Errorf("error getmessage from queue: %v", err)
			return
		}

		if len(retResponse.Messages) > 0 {

			processedReceiptHandles := make([]*sqs.DeleteMessageBatchRequestEntry, len(retResponse.Messages))

			for i, mess := range retResponse.Messages {
				revel.AppLog.Infof("message =[%v] == %+v", i, mess.String())

				err = doSendAft(mess.Body)
				if err != nil {
					revel.AppLog.Errorf("could not doSendAft: %v", err)
					return
				}

				processedReceiptHandles[i] = &sqs.DeleteMessageBatchRequestEntry{
					Id:            mess.MessageId,
					ReceiptHandle: mess.ReceiptHandle,
				}
			}

			err = sqsConn.DeleteBatchFromQueue(&sqs.DeleteMessageBatchInput{
				QueueUrl: theQueue.QueueUrl,
				Entries:  processedReceiptHandles,
			})
			if err != nil {
				revel.AppLog.Errorf("failed to delete messages batch from queue: %+v", err)
				return
			}

		}

		if len(retResponse.Messages) == 0 {
			revel.AppLog.Infof("processATSendRequests:: No sqs messages in at requests, sleeping for 10")
			time.Sleep(time.Second * 10)
		}
	}
}

func doSendAft(data *string) error {
	theData := map[string]interface{}{}
	err := json.Unmarshal([]byte(*data), &theData)
	if err != nil {
		return err
	}

	recs := make([]*entities.SMSRecipient, 0)
	for _, rec := range theData["recs"].([]interface{}) {
		recs = append(recs, &entities.SMSRecipient{
			Phone: rec.(map[string]interface{})["phone"].(string),
		})
	}
	_, err = africasTalkingSender.Send(&entities.SendRequest{
		SenderID: theData["sender_id"].(string),
		Message:  theData["message"].(string),
		To:       recs,
	})
	if err != nil {
		return err
	}
	return nil
}

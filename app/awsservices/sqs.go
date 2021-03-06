package awsservices

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSClient holds a connection to AWS ECR.
type SQSClient struct {
	client *sqs.SQS
	region string
}

// NewSQSClient initializes an SQSClient.
func NewSQSClient() *SQSClient {
	region := "eu-west-1"
	return &SQSClient{
		client: sqs.New(session.New(), aws.NewConfig().WithRegion(region)),
		region: region,
	}
}

// GetQueueURL returns the url of sqs queue
func (c *SQSClient) GetQueueURL(
	queue string,
) (*sqs.GetQueueUrlOutput, error) {
	result, err := c.client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *SQSClient) MessageToQueue(
	queueURL *string,
	data map[string]interface{},
) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = c.client.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"User": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("Eutychus Towett"),
			},
		},
		MessageBody: aws.String(string(jsonData)),
		QueueUrl:    queueURL,
	})
	return err
}

func (c *SQSClient) GetMessageFromQueue(
	queueURL *string,
) (*sqs.ReceiveMessageOutput, error) {
	messageResponse, err := c.client.ReceiveMessage(&sqs.ReceiveMessageInput{QueueUrl: queueURL})

	if err != nil {
		return messageResponse, err
	}
	return messageResponse, err
}

func (c *SQSClient) DeleteBatchFromQueue(
	deleteMessageRequest *sqs.DeleteMessageBatchInput,
) error {
	_, err := c.client.DeleteMessageBatch(deleteMessageRequest)

	if err != nil {
		return err
	}
	return nil
}

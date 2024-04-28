package repository

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type BasicSQSClient interface {
	GetQueueAttributes(queueUrl string) (*sqs.GetQueueAttributesOutput, error)
	GetQueueURL(queueName string) (string, error)
	ReceiveMessage(queueUrl string, maxMessages *int32, waitTime *int32) (*sqs.ReceiveMessageOutput, error)
	SendMessage(queueUrl string, messageGroupID string, messageBody string) (*sqs.SendMessageOutput, error)
	SendMessageBatch(queueUrl string, messages []types.SendMessageBatchRequestEntry) (*sqs.SendMessageBatchOutput, error)
	DeleteMessage(queueUrl string, message *types.Message) (*sqs.DeleteMessageOutput, error)
	DeleteMessageBatch(queueUrl string, messages []types.Message) (*sqs.DeleteMessageBatchOutput, error)
}

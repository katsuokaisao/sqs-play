package repository

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/katsuokaisao/sql-play/domain"
)

type ExampleSQSClient interface {
	GetQueueURL() (string, error)
	SetQueueURL(queueURL string)
	GetQueueAttributes() (map[string]string, error)
	SendMessage(message *domain.Example) (string, error)
	SendMessageBatch(messages []domain.Example) (*sqs.SendMessageBatchOutput, error)
	ReceiveMessage() ([]domain.ExampleMsg, error)
	DeleteMessage(msg *domain.ExampleMsg) error
	DeleteMessageBatch(msgs []domain.ExampleMsg) (*sqs.DeleteMessageBatchOutput, error)
}

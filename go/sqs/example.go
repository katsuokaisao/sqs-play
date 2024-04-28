package sqs

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/katsuokaisao/sql-play/domain"
	"github.com/katsuokaisao/sql-play/repository"
)

type exampleSQSClient struct {
	cli            repository.BasicSQSClient
	queueName      string
	queueURL       string
	messageGroupID string
	maxMessages    int32
	waitTime       int32
}

func NewExampleSQSClient(
	cli repository.BasicSQSClient,
	queueName string,
	queueURL string,
	maxMessages int32,
	waitTime int32,
) repository.ExampleSQSClient {
	return &exampleSQSClient{
		cli:            cli,
		queueName:      queueName,
		queueURL:       queueURL,
		messageGroupID: "example",
		maxMessages:    maxMessages,
		waitTime:       waitTime,
	}
}

func (c *exampleSQSClient) GetQueueURL() (string, error) {
	if c.queueURL != "" {
		return c.queueURL, nil
	}

	queueURL, err := c.cli.GetQueueURL(c.queueName)
	if err != nil {
		return "", fmt.Errorf("failed to get queue url: %w", err)
	}
	return queueURL, nil
}

func (c *exampleSQSClient) SetQueueURL(queueURL string) {
	c.queueURL = queueURL
}

func (c *exampleSQSClient) GetQueueAttributes() (map[string]string, error) {
	attribute, err := c.cli.GetQueueAttributes(c.queueURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue attributes: %w", err)
	}

	return attribute.Attributes, nil
}

func (c *exampleSQSClient) SendMessage(message *domain.Example) (string, error) {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("failed to marshal message: %w", err)
	}

	output, err := c.cli.SendMessage(c.queueURL, c.messageGroupID, string(messageJSON))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	return *output.MessageId, nil
}

func (c *exampleSQSClient) SendMessageBatch(messages []domain.Example) (*sqs.SendMessageBatchOutput, error) {
	var entries []types.SendMessageBatchRequestEntry
	for i := range messages {
		message := messages[i]
		messageJSON, err := json.Marshal(message)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal message: %w", err)
		}
		entries = append(entries, types.SendMessageBatchRequestEntry{
			Id:             aws.String(fmt.Sprintf("%d", i)),
			MessageBody:    aws.String(string(messageJSON)),
			MessageGroupId: aws.String(c.messageGroupID),
		})
	}

	output, err := c.cli.SendMessageBatch(c.queueURL, entries)
	if err != nil {
		return output, fmt.Errorf("failed to send message batch: %w", err)
	}

	return output, nil
}

func (c *exampleSQSClient) ReceiveMessage() ([]domain.ExampleMsg, error) {
	result, err := c.cli.ReceiveMessage(c.queueURL, &c.maxMessages, &c.waitTime)
	if err != nil {
		return nil, fmt.Errorf("failed to receive message: %w", err)
	}

	messages := make([]domain.ExampleMsg, 0, len(result.Messages))
	for i := range result.Messages {
		msg := result.Messages[i]
		var data domain.Example
		if err := json.Unmarshal([]byte(*msg.Body), &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %w", err)
		}

		messages = append(messages, domain.ExampleMsg{
			Data: data,
			Msg:  &msg,
		})
	}

	return messages, nil
}

func (c *exampleSQSClient) DeleteMessage(msg *domain.ExampleMsg) error {
	if msg.Msg == nil {
		return errors.New("message is nil")
	}

	_, err := c.cli.DeleteMessage(c.queueURL, msg.Msg)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (c *exampleSQSClient) DeleteMessageBatch(msgs []domain.ExampleMsg) (*sqs.DeleteMessageBatchOutput, error) {
	if len(msgs) == 0 {
		return nil, errors.New("messages is empty")
	}

	var messages []types.Message
	for i := range msgs {
		msg := msgs[i]
		if msg.Msg == nil {
			return nil, errors.New("message is nil")
		}
		messages = append(messages, *msg.Msg)
	}

	output, err := c.cli.DeleteMessageBatch(c.queueURL, messages)
	if err != nil {
		return output, fmt.Errorf("failed to delete message batch: %w", err)
	}

	return output, nil
}

package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"github.com/katsuokaisao/sql-play/repository"
)

type basicSQSClient struct {
	cli *sqs.Client
}

func NewBasicSQSClient(cli *sqs.Client) repository.BasicSQSClient {
	return &basicSQSClient{cli: cli}
}

func (actor *basicSQSClient) GetQueueAttributes(queueUrl string) (*sqs.GetQueueAttributesOutput, error) {
	arnAttributeName := types.QueueAttributeNameQueueArn
	attribute, err := actor.cli.GetQueueAttributes(context.TODO(), &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueUrl),
		AttributeNames: []types.QueueAttributeName{arnAttributeName},
	})
	return attribute, err
}

func (actor *basicSQSClient) GetQueueURL(queueName string) (string, error) {
	queueUrl, err := actor.cli.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", err
	}
	if queueUrl.QueueUrl == nil {
		return "", fmt.Errorf("queue url is empty")
	}

	return *queueUrl.QueueUrl, nil
}

func (actor *basicSQSClient) ReceiveMessage(queueUrl string, maxMessages *int32, waitTime *int32) (*sqs.ReceiveMessageOutput, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueUrl),
	}
	if maxMessages != nil {
		input.MaxNumberOfMessages = *maxMessages
	}
	if waitTime != nil {
		input.WaitTimeSeconds = *waitTime
	}

	result, err := actor.cli.ReceiveMessage(context.TODO(), input)
	return result, err
}

func (actor *basicSQSClient) SendMessage(queueUrl string, messageGroupID string, messageBody string) (*sqs.SendMessageOutput, error) {
	output, err := actor.cli.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:       aws.String(queueUrl),
		MessageGroupId: aws.String(messageGroupID),
		MessageBody:    aws.String(messageBody),
	})
	return output, err
}

func (actor *basicSQSClient) SendMessageBatch(queueUrl string, messages []types.SendMessageBatchRequestEntry) (*sqs.SendMessageBatchOutput, error) {
	output, err := actor.cli.SendMessageBatch(context.TODO(), &sqs.SendMessageBatchInput{
		QueueUrl: aws.String(queueUrl),
		Entries:  messages,
	})
	return output, err
}

func (actor *basicSQSClient) DeleteMessage(queueUrl string, message *types.Message) (*sqs.DeleteMessageOutput, error) {
	output, err := actor.cli.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: message.ReceiptHandle,
	})
	return output, err
}

func (actor *basicSQSClient) DeleteMessageBatch(queueUrl string, messages []types.Message) (*sqs.DeleteMessageBatchOutput, error) {
	entries := make([]types.DeleteMessageBatchRequestEntry, len(messages))
	for msgIndex := range messages {
		entries[msgIndex].Id = aws.String(uuid.New().String())
		entries[msgIndex].ReceiptHandle = messages[msgIndex].ReceiptHandle
	}
	output, err := actor.cli.DeleteMessageBatch(context.TODO(), &sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(queueUrl),
	})
	return output, err
}

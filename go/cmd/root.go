package cmd

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/katsuokaisao/sql-play/domain"
	"github.com/katsuokaisao/sql-play/repository"
	"github.com/katsuokaisao/sql-play/sqs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sendMessageBatchCmd)
	rootCmd.AddCommand(sendMessageCmd)
	rootCmd.AddCommand(receiveMessageCmd)
	rootCmd.AddCommand(getQueueAttributeCmd)
	rootCmd.AddCommand(deleteMessageCmd)
	rootCmd.AddCommand(deleteMessageBatchCmd)
}

var rootCmd = &cobra.Command{
	Use: "",
}

func Execute() error {
	return rootCmd.Execute()
}

func initSQS() repository.ExampleSQSClient {
	e := domain.Env{
		SQS: &domain.SQSEnv{},
	}
	if err := env.Parse(&e); err != nil {
		panic(fmt.Errorf("failed to parse env: %w", err))
	}
	if e.SQS == nil || (e.SQS.QueueName == "" && e.SQS.QueueURL == "") {
		panic("SQS is nil")
	}

	awsCfg, err := sqs.ProvideAWSConfig()
	if err != nil {
		panic(fmt.Errorf("failed to provide aws config: %w", err))
	}
	sqsClient := sqs.ProviderSQSClient(awsCfg)
	sqsBasicClient := sqs.NewBasicSQSClient(sqsClient)
	exampleSQSClient := sqs.NewExampleSQSClient(
		sqsBasicClient,
		e.SQS.QueueName,
		e.SQS.QueueURL,
		int32(e.SQS.MaxNum),
		int32(e.SQS.WaitTime),
	)

	queueURL, err := exampleSQSClient.GetQueueURL()
	if err != nil {
		panic(fmt.Errorf("failed to get queue url: %w", err))
	}
	if queueURL == "" {
		panic("queueURL is empty")
	}
	fmt.Printf("QueueURL: %s\n", queueURL)
	exampleSQSClient.SetQueueURL(queueURL)

	return exampleSQSClient
}

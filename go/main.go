package main

import (
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/caarlos0/env"
	"github.com/katsuokaisao/sql-play/domain"
	"github.com/katsuokaisao/sql-play/sqs"
)

func main() {
	e := domain.Env{
		SQS: &domain.SQSEnv{},
	}
	if err := env.Parse(&e); err != nil {
		panic(fmt.Errorf("failed to parse env: %w", err))
	}
	if e.SQS == nil || (e.SQS.QueueName == "" && e.SQS.QueueURL == "") {
		panic("SQS is nil")
	}
	fmt.Printf("SQS: %v\n", e.SQS)

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

	queueAttrs, err := exampleSQSClient.GetQueueAttributes()
	if err != nil {
		panic(fmt.Errorf("failed to get queue attributes: %w", err))
	}
	fmt.Printf("queueAttrs: %v\n", queueAttrs)

	{
		for i := 0; i < 10; i++ {
			message, err := exampleSQSClient.ReceiveMessage()
			if err != nil {
				panic(fmt.Errorf("failed to receive message: %w", err))
			}
			fmt.Printf("before len(message): %v\n", len(message))

			if len(message) == 0 {
				break
			}

			for i := range message {
				msg := message[i]
				fmt.Printf("success ReceiveMessage message: %v\n", msg.Data)
				if err = exampleSQSClient.DeleteMessage(&msg); err != nil {
					panic(err)
				}
				fmt.Printf("success DeleteMessage messageID: %s\n", *msg.Msg.MessageId)
			}
		}
	}

	{
		msg := domain.NewExample()
		messageID, err := exampleSQSClient.SendMessage(msg)
		if err != nil {
			panic(fmt.Errorf("failed to send message: %w", err))
		}
		fmt.Printf("success SendMessage messageID: %s\n", messageID)

		time.Sleep(1 * time.Second)

		message, err := exampleSQSClient.ReceiveMessage()
		if err != nil {
			panic(fmt.Errorf("failed to receive message: %w", err))
		}
		fmt.Printf("ReceiveMessage len(message): %v\n", len(message))

		for i := range message {
			msg := message[i]
			fmt.Printf("ReceiveMessage message: %v\n", msg.Data)
			if err = exampleSQSClient.DeleteMessage(&msg); err != nil {
				panic(err)
			}
			fmt.Printf("DeleteMessage messageID: %s\n", *msg.Msg.MessageId)
		}
	}

	fmt.Println("--------------------------------")

	{
		msgs := make([]domain.Example, 0, 1024)
		for i := 0; i < 10; i++ {
			msgs = append(msgs, domain.Example{
				Filed1: fmt.Sprintf("Filed1-%d", i),
				Filed2: i,
				Filed3: i%2 == 0,
				Field4: float64(i),
				Field5: []string{fmt.Sprintf("Field5-%d", i), fmt.Sprintf("Field5-%d", i+1)},
				Field6: map[string]string{
					"key1": fmt.Sprintf("value1-%d", i),
					"key2": fmt.Sprintf("value2-%d", i),
				},
				Field7: pointer.ToString(fmt.Sprintf("Field7-%d", i)),
			})
		}
		output, err := exampleSQSClient.SendMessageBatch(msgs)
		if err != nil {
			panic(fmt.Errorf("failed to send message batch: %w", err))
		}
		for i := range output.Failed {
			f := output.Failed[i]
			if f.Id != nil && f.Code != nil && f.Message != nil {
				fmt.Printf("id: %s, code: %s, msg: %s, sender_fault: %v \n", *f.Id, *f.Code, *f.Message, f.SenderFault)
			} else {
				fmt.Printf("failed message: %v\n", f)
			}
		}
		for i := range output.Successful {
			s := output.Successful[i]
			fmt.Printf("success SendMessageBatch messageID: %v\n", *s.MessageId)
		}

		time.Sleep(3 * time.Second)

		message2, err := exampleSQSClient.ReceiveMessage()
		if err != nil {
			panic(fmt.Errorf("failed to receive message: %w", err))
		}
		if len(message2) != 0 {
			for _, m := range message2 {
				fmt.Printf("success ReceiveMessage message: %v\n", m.Data)
			}

			deleteMessageBatchOutput, err := exampleSQSClient.DeleteMessageBatch(message2)
			if err != nil {
				panic(fmt.Errorf("failed to delete message batch: %w", err))
			}
			for _, f := range deleteMessageBatchOutput.Failed {
				if f.Id != nil && f.Code != nil && f.Message != nil {
					fmt.Printf("failed id: %s, code: %s, msg: %s, sender_fault: %v \n", *f.Id, *f.Code, *f.Message, f.SenderFault)
				} else {
					fmt.Printf("failed: %v\n", f)
				}
			}
			for _, s := range deleteMessageBatchOutput.Successful {
				fmt.Printf("success DeleteMessageBatch messageID: %v\n", *s.Id)
			}
		}

	}

	fmt.Println("done")
}

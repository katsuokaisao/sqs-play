package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteMessageBatchCmd = &cobra.Command{
	Use: "deleteMessageBatch",
	Run: func(cmd *cobra.Command, args []string) {
		exampleSQSClient := initSQS()
		message, err := exampleSQSClient.ReceiveMessage()
		if err != nil {
			panic(fmt.Errorf("failed to receive message: %w", err))
		}

		if len(message) == 0 {
			fmt.Println("no message")
			return
		}

		for i := range message {
			msg := message[i]
			fmt.Printf("success ReceiveMessage message: %v\n", msg.Data)
		}

		deleteMessageBatchOutput, err := exampleSQSClient.DeleteMessageBatch(message)
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

		fmt.Println("done")
	},
}

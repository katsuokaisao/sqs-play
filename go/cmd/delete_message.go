package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteMessageCmd = &cobra.Command{
	Use: "deleteMessage",
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
			if err = exampleSQSClient.DeleteMessage(&msg); err != nil {
				panic(fmt.Errorf("failed to delete message: %w", err))
			}
			fmt.Printf("success DeleteMessage messageID: %s\n", *msg.Msg.MessageId)
		}
	},
}

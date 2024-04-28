package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var receiveMessageCmd = &cobra.Command{
	Use: "receiveMessage",
	Run: func(cmd *cobra.Command, args []string) {
		exampleSQSClient := initSQS()
		message, err := exampleSQSClient.ReceiveMessage()
		if err != nil {
			panic(fmt.Errorf("failed to receive message: %w", err))
		}

		for i := range message {
			msg := message[i]
			fmt.Printf("success ReceiveMessage message: %v\n", msg.Data)
		}
	},
}

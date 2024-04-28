package cmd

import (
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/katsuokaisao/sql-play/domain"
	"github.com/spf13/cobra"
)

var sendMessageCmd = &cobra.Command{
	Use: "sendMessage",
	Run: func(cmd *cobra.Command, args []string) {
		exampleSQSClient := initSQS()
		msg := domain.Example{
			Filed1: "Filed1",
			Filed2: 1,
			Filed3: true,
			Field4: 1.1,
			Field5: []string{"Field5-1", "Field5-2"},
			Field6: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			Field7: pointer.ToString("Field7"),
		}

		messageID, err := exampleSQSClient.SendMessage(&msg)
		if err != nil {
			panic(fmt.Errorf("failed to send message: %w", err))
		}
		fmt.Printf("success SendMessage messageID: %v\n", messageID)

		fmt.Println("done")
	},
}

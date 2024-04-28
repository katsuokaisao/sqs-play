package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getQueueAttributeCmd = &cobra.Command{
	Use: "getQueueAttribute",
	Run: func(cmd *cobra.Command, args []string) {
		exampleSQSClient := initSQS()
		queueAttrs, err := exampleSQSClient.GetQueueAttributes()
		if err != nil {
			panic(fmt.Errorf("failed to get queue attributes: %w", err))
		}
		fmt.Printf("queueAttrs: %v\n", queueAttrs)
	},
}

package cmd

import (
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/katsuokaisao/sql-play/domain"
	"github.com/spf13/cobra"
)

var sendMessageBatchCmd = &cobra.Command{
	Use: "sendMessageBatch",
	Run: func(cmd *cobra.Command, args []string) {
		exampleSQSClient := initSQS()
		msgs := make([]domain.Example, 0, 1024)
		for i := 0; i < 5; i++ {
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

		fmt.Println("done")
	},
}

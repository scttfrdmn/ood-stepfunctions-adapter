package cmd

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/ood-stepfunctions-adapter/internal/stepfunc"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <execution-arn>",
	Short: "Stop a Step Functions execution",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := stepfunc.New(ctx, region)
		if err != nil {
			return err
		}
		if err := client.StopExecution(ctx, args[0], "Cancelled via OOD"); err != nil {
			return err
		}
		fmt.Printf("Execution %s stopped\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

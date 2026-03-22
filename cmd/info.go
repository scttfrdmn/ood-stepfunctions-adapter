package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/scttfrdmn/ood-stepfunctions-adapter/internal/stepfunc"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <execution-arn>",
	Short: "Print full Step Functions execution details as JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := stepfunc.New(ctx, region)
		if err != nil {
			return err
		}
		detail, err := client.DescribeExecution(ctx, args[0])
		if err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(detail)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

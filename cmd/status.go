package cmd

import (
	"context"
	"encoding/json"
	"os"

	internalood "github.com/scttfrdmn/ood-stepfunctions-adapter/internal/ood"
	"github.com/scttfrdmn/ood-stepfunctions-adapter/internal/stepfunc"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <execution-arn>",
	Short: "Get the status of a Step Functions execution",
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

		js := internalood.JobStatus{
			ID:     args[0],
			Status: stepfunc.SfnStateToOod(detail.Status),
		}
		if detail.Cause != nil {
			js.Message = *detail.Cause
		}

		return json.NewEncoder(os.Stdout).Encode(js)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

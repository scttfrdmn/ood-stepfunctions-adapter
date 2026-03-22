package cmd

import (
	"github.com/spf13/cobra"
)

var (
	region          string
	stateMachineArn string
)

var rootCmd = &cobra.Command{
	Use:   "ood-stepfunctions-adapter",
	Short: "OOD compute adapter for AWS Step Functions",
	Long:  "Translates Open OnDemand job submissions to AWS Step Functions API calls.",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "AWS region")
	rootCmd.PersistentFlags().StringVar(&stateMachineArn, "state-machine-arn", "", "ARN of the Step Functions state machine")
}

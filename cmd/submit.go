package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/scttfrdmn/ood-stepfunctions-adapter/internal/stepfunc"
	"github.com/spf13/cobra"
)

// JobSpec is the Step Functions-specific job submission payload.
type JobSpec struct {
	StateMachineArn string                 `json:"state_machine_arn,omitempty"`
	Input           map[string]interface{} `json:"input,omitempty"`
	JobName         string                 `json:"job_name,omitempty"`
}

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit an OOD job to AWS Step Functions",
	Long:  "Reads a JSON job spec from stdin and starts a Step Functions execution.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var spec JobSpec
		if err := json.NewDecoder(os.Stdin).Decode(&spec); err != nil {
			return fmt.Errorf("decode job spec: %w", err)
		}

		effectiveSMA := stateMachineArn
		if effectiveSMA == "" {
			effectiveSMA = spec.StateMachineArn
		}
		if effectiveSMA == "" {
			return fmt.Errorf("--state-machine-arn is required (or set state_machine_arn in job spec)")
		}

		ctx := context.Background()
		client, err := stepfunc.New(ctx, region)
		if err != nil {
			return err
		}

		executionArn, err := client.StartExecution(ctx, stepfunc.ExecutionSpec{
			StateMachineArn: effectiveSMA,
			Input:           spec.Input,
			JobName:         spec.JobName,
		})
		if err != nil {
			return err
		}

		fmt.Println(executionArn)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}

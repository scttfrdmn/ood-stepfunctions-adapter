// Package stepfunc wraps the AWS Step Functions API for the OOD adapter.
package stepfunc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awssfn "github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

// Client wraps the AWS Step Functions client.
type Client struct {
	svc    *awssfn.Client
	region string
}

// New creates a Step Functions client using the default AWS credential chain.
func New(ctx context.Context, region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}
	return &Client{svc: awssfn.NewFromConfig(cfg), region: region}, nil
}

// ExecutionSpec holds the parameters for a Step Functions execution.
type ExecutionSpec struct {
	StateMachineArn string
	Input           map[string]interface{}
	JobName         string
}

// StartExecution starts a Step Functions execution and returns the execution ARN.
func (c *Client) StartExecution(ctx context.Context, spec ExecutionSpec) (string, error) {
	inputJSON := "{}"
	if len(spec.Input) > 0 {
		b, err := json.Marshal(spec.Input)
		if err != nil {
			return "", fmt.Errorf("marshal execution input: %w", err)
		}
		inputJSON = string(b)
	}

	input := &awssfn.StartExecutionInput{
		StateMachineArn: aws.String(spec.StateMachineArn),
		Input:           aws.String(inputJSON),
	}
	if spec.JobName != "" {
		input.Name = aws.String(spec.JobName)
	}

	out, err := c.svc.StartExecution(ctx, input)
	if err != nil {
		return "", fmt.Errorf("sfn StartExecution: %w", err)
	}
	return aws.ToString(out.ExecutionArn), nil
}

// DescribeExecution returns the current detail of a Step Functions execution.
func (c *Client) DescribeExecution(ctx context.Context, executionArn string) (*awssfn.DescribeExecutionOutput, error) {
	out, err := c.svc.DescribeExecution(ctx, &awssfn.DescribeExecutionInput{
		ExecutionArn: aws.String(executionArn),
	})
	if err != nil {
		return nil, fmt.Errorf("sfn DescribeExecution: %w", err)
	}
	return out, nil
}

// StopExecution stops a Step Functions execution.
func (c *Client) StopExecution(ctx context.Context, executionArn, cause string) error {
	_, err := c.svc.StopExecution(ctx, &awssfn.StopExecutionInput{
		ExecutionArn: aws.String(executionArn),
		Cause:        aws.String(cause),
	})
	if err != nil {
		return fmt.Errorf("sfn StopExecution: %w", err)
	}
	return nil
}

// SfnStateToOod maps a Step Functions execution status to an OOD status string.
func SfnStateToOod(s types.ExecutionStatus) string {
	switch s {
	case types.ExecutionStatusRunning:
		return "running"
	case types.ExecutionStatusSucceeded:
		return "completed"
	case types.ExecutionStatusFailed, types.ExecutionStatusTimedOut:
		return "failed"
	case types.ExecutionStatusAborted:
		return "cancelled"
	default:
		return "undetermined"
	}
}

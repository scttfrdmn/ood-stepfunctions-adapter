//go:build integration

package stepfunc_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/scttfrdmn/ood-stepfunctions-adapter/internal/stepfunc"
	substrate "github.com/scttfrdmn/substrate"
)

// createTestStateMachine registers a minimal Pass state machine in substrate.
// Returns the state machine ARN.
func createTestStateMachine(t *testing.T, ctx context.Context, endpointURL string) string {
	t.Helper()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithBaseEndpoint(endpointURL),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		t.Fatalf("createTestStateMachine: load config: %v", err)
	}
	sfnSvc := sfn.NewFromConfig(cfg)
	out, err := sfnSvc.CreateStateMachine(ctx, &sfn.CreateStateMachineInput{
		Name: aws.String("ood-test"),
		Type: sfntypes.StateMachineTypeStandard,
		Definition: aws.String(`{"Comment":"test","StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}}`),
		RoleArn: aws.String("arn:aws:iam::123456789012:role/test-role"),
	})
	if err != nil {
		t.Fatalf("createTestStateMachine: CreateStateMachine: %v", err)
	}
	return aws.ToString(out.StateMachineArn)
}

func TestStartExecution_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()

	// Create the state machine before starting an execution.
	smArn := createTestStateMachine(t, ctx, ts.URL)
	t.Logf("state machine ARN: %s", smArn)

	client, err := stepfunc.New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	spec := stepfunc.ExecutionSpec{
		StateMachineArn: smArn,
		Input:           map[string]interface{}{"key": "value"},
		JobName:         "test-execution",
	}

	executionArn, err := client.StartExecution(ctx, spec)
	if err != nil {
		t.Fatalf("StartExecution: %v", err)
	}
	if executionArn == "" {
		t.Fatal("expected non-empty execution ARN")
	}
	t.Logf("execution ARN: %s", executionArn)

	// Describe the execution and verify status.
	detail, err := client.DescribeExecution(ctx, executionArn)
	if err != nil {
		t.Fatalf("DescribeExecution: %v", err)
	}
	if aws.ToString(detail.ExecutionArn) == "" {
		t.Error("expected non-empty ExecutionArn in response")
	}
	t.Logf("execution status: %s", detail.Status)
	if detail.Status == "" {
		t.Error("expected non-empty Status")
	}

	// Stop the execution.
	err = client.StopExecution(ctx, executionArn, "test teardown")
	if err != nil {
		t.Fatalf("StopExecution: %v", err)
	}
	t.Logf("execution stopped")
}

func TestDescribeExecution_NotFound_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()
	client, err := stepfunc.New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = client.DescribeExecution(ctx, "arn:aws:states:us-east-1:123456789012:execution:does-not-exist:no-such-exec")
	if err == nil {
		t.Fatal("expected error for non-existent execution, got nil")
	}
	t.Logf("error (expected): %v", err)
}

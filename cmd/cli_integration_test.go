//go:build integration

package cmd

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awssfn "github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
	substrate "github.com/scttfrdmn/substrate"
)

// runSFNCmd pipes stdinData into os.Stdin, captures os.Stdout, executes
// rootCmd with args, and returns the trimmed output.
func runSFNCmd(t *testing.T, stdinData string, args ...string) string {
	t.Helper()

	if stdinData != "" {
		stdinR, stdinW, err := os.Pipe()
		if err != nil {
			t.Fatalf("create stdin pipe: %v", err)
		}
		origStdin := os.Stdin
		os.Stdin = stdinR
		defer func() {
			os.Stdin = origStdin
			stdinR.Close()
		}()
		if _, err := io.WriteString(stdinW, stdinData); err != nil {
			t.Fatalf("write stdin pipe: %v", err)
		}
		stdinW.Close()
	}

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	origStdout := os.Stdout
	os.Stdout = stdoutW
	defer func() { os.Stdout = origStdout }()

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.SetArgs(args)
	execErr := rootCmd.Execute()

	stdoutW.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	if _, readErr := io.Copy(&buf, stdoutR); readErr != nil {
		t.Fatalf("read stdout pipe: %v", readErr)
	}
	stdoutR.Close()

	if execErr != nil {
		t.Fatalf("rootCmd.Execute(%v): %v", args, execErr)
	}
	return strings.TrimSpace(buf.String())
}

func TestCLISubmitStatusDelete_StepFunctions_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithBaseEndpoint(ts.URL),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		t.Fatalf("load aws config: %v", err)
	}
	raw := awssfn.NewFromConfig(cfg)

	// Create the state machine that the CLI will target.
	smOut, err := raw.CreateStateMachine(ctx, &awssfn.CreateStateMachineInput{
		Name:       aws.String("ood-cli-sfn"),
		Type:       sfntypes.StateMachineTypeStandard,
		Definition: aws.String(`{"Comment":"test","StartAt":"Pass","States":{"Pass":{"Type":"Pass","End":true}}}`),
		RoleArn:    aws.String("arn:aws:iam::123456789012:role/SFNRole"),
	})
	if err != nil {
		t.Fatalf("CreateStateMachine: %v", err)
	}
	smArn := aws.ToString(smOut.StateMachineArn)
	t.Logf("created state machine: %s", smArn)

	// submit
	spec := `{"job_name":"cli-sfn-test","input":{"key":"value"}}`
	executionArn := runSFNCmd(t, spec,
		"submit",
		"--region", "us-east-1",
		"--state-machine-arn", smArn,
	)
	if executionArn == "" {
		t.Fatal("submit: expected non-empty execution ARN")
	}
	if !strings.Contains(executionArn, "arn:aws") {
		t.Errorf("submit: execution ARN does not look like an ARN: %s", executionArn)
	}
	t.Logf("submitted execution: %s", executionArn)

	// status
	statusOut := runSFNCmd(t, "",
		"status",
		"--region", "us-east-1",
		executionArn,
	)
	t.Logf("status output: %s", statusOut)
	if !strings.Contains(statusOut, "running") && !strings.Contains(statusOut, "completed") {
		t.Errorf("status output does not contain a recognised status: %s", statusOut)
	}

	// delete
	deleteOut := runSFNCmd(t, "",
		"delete",
		"--region", "us-east-1",
		executionArn,
	)
	t.Logf("delete output: %s", deleteOut)
	if !strings.Contains(deleteOut, executionArn) {
		t.Errorf("delete output does not reference execution ARN %q: %s", executionArn, deleteOut)
	}
}

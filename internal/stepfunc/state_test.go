package stepfunc

import (
	"testing"

	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

func TestSfnStateToOod(t *testing.T) {
	tests := []struct {
		state sfntypes.ExecutionStatus
		want  string
	}{
		{sfntypes.ExecutionStatusRunning, "running"},
		{sfntypes.ExecutionStatusSucceeded, "completed"},
		{sfntypes.ExecutionStatusFailed, "failed"},
		{sfntypes.ExecutionStatusTimedOut, "failed"},
		{sfntypes.ExecutionStatusAborted, "cancelled"},
		{"UNKNOWN_STATE", "undetermined"},
	}
	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			got := SfnStateToOod(tt.state)
			if got != tt.want {
				t.Errorf("SfnStateToOod(%q) = %q, want %q", tt.state, got, tt.want)
			}
		})
	}
}

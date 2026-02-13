package cli

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

type MockCommandExecutor struct {
	RunFunc func(name string, args []string, env []string, stdout, stderr io.Writer) error
}

func (m *MockCommandExecutor) RunCommand(name string, args []string, env []string, stdout, stderr io.Writer) error {
	return m.RunFunc(name, args, env, stdout, stderr)
}

func TestCopilotCLI_Install(t *testing.T) {
	callCount := 0
	mock := &MockCommandExecutor{
		RunFunc: func(name string, args []string, env []string, stdout, stderr io.Writer) error {
			callCount++
			return nil
		},
	}

	copilot := &CopilotCLI{}
	err := copilot.Install(mock)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedCalls := 6 // 5 for GH CLI + 1 for Extension
	if callCount != expectedCalls {
		t.Errorf("expected %d command calls, got %d", expectedCalls, callCount)
	}
}

func TestCopilotCLI_Auth(t *testing.T) {
	copilot := &CopilotCLI{}
	err := copilot.Auth(nil, "token")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if copilot.Token != "token" {
		t.Errorf("expected token to be set")
	}
}

func TestCopilotCLI_Run(t *testing.T) {
	mock := &MockCommandExecutor{
		RunFunc: func(name string, args []string, env []string, stdout, stderr io.Writer) error {
			if name != "gh" {
				return fmt.Errorf("expected gh command")
			}
			foundToken := false
			for _, e := range env {
				if strings.HasPrefix(e, "COPILOT_GITHUB_TOKEN=") {
					foundToken = true
				}
			}
			if !foundToken {
				return fmt.Errorf("expected COPILOT_GITHUB_TOKEN in env")
			}
			return nil
		},
	}

	copilot := &CopilotCLI{Token: "test-token"} // Manually set token or use Auth
	err := copilot.Run(mock, "prompt", "model")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCopilotCLI_Run_NoToken(t *testing.T) {
	copilot := &CopilotCLI{}
	err := copilot.Run(nil, "prompt", "model")
	if err == nil {
		t.Error("expected error when token is missing")
	}
}

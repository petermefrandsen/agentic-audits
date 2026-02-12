package main

import (
	"fmt"
	"io"
	"testing"
)

type MockCommandExecutor struct {
	RunFunc func(name string, args []string, env []string, stdout, stderr io.Writer) error
}

func (m *MockCommandExecutor) RunCommand(name string, args []string, env []string, stdout, stderr io.Writer) error {
	return m.RunFunc(name, args, env, stdout, stderr)
}

func TestInstallGitHubCLI(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		callCount := 0
		mock := &MockCommandExecutor{
			RunFunc: func(name string, args []string, env []string, stdout, stderr io.Writer) error {
				callCount++
				return nil
			},
		}
		err := installGitHubCLI(mock)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if callCount != 5 {
			t.Errorf("expected 5 command calls, got %d", callCount)
		}
	})

	t.Run("Failure", func(t *testing.T) {
		mock := &MockCommandExecutor{
			RunFunc: func(name string, args []string, env []string, stdout, stderr io.Writer) error {
				return fmt.Errorf("fail")
			},
		}
		err := installGitHubCLI(mock)
		if err == nil {
			t.Error("expected error on command failure")
		}
	})
}

func TestInstallCopilotExtension(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock := &MockCommandExecutor{
			RunFunc: func(name string, args []string, env []string, stdout, stderr io.Writer) error {
				return nil
			},
		}
		err := installCopilotExtension(mock)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestContains(t *testing.T) {
	if !contains("hello world", "world") {
		t.Error("expected true")
	}
	if contains("hello", "world") {
		t.Error("expected false")
	}
}

package main

import (
	"io"
	"testing"
)

type MockCommandExecutor struct {
	RunFunc func(name string, args []string, env []string, stdout, stderr io.Writer) error
}

func (m *MockCommandExecutor) RunCommand(name string, args []string, env []string, stdout, stderr io.Writer) error {
	return m.RunFunc(name, args, env, stdout, stderr)
}

// Tests for installGitHubCLI and installCopilotExtension are moved or should be moved to cli/copilot_test.go

func TestContains(t *testing.T) {
	if !contains("hello world", "world") {
		t.Error("expected true")
	}
	if contains("hello", "world") {
		t.Error("expected false")
	}
}

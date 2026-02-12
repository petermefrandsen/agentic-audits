package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)


func TestConstructFullPrompt(t *testing.T) {
	os.Setenv("GITHUB_REPOSITORY", "test/repo")
	defer os.Unsetenv("GITHUB_REPOSITORY")

	mission := "test mission"
	opts := AgentOptions{
		ContextFiles: ".",
		DryRun:       false,
	}
	webSources := "web info"

	prompt := constructFullPrompt(mission, opts, webSources)

	if !strings.Contains(prompt, "test mission") {
		t.Error("prompt should contain mission")
	}
	if !strings.Contains(prompt, "web info") {
		t.Error("prompt should contain web sources")
	}
	if !strings.Contains(prompt, "MANDATORY: Pull Request Creation") {
		t.Error("prompt should contain PR instructions when not dry-run")
	}

	// Test dry-run
	opts.DryRun = true
	prompt = constructFullPrompt(mission, opts, webSources)
	if strings.Contains(prompt, "MANDATORY: Pull Request Creation") {
		t.Error("prompt should NOT contain PR instructions when dry-run")
	}
	if !strings.Contains(prompt, "dry_run is set to TRUE") {
		t.Error("prompt should contain dry-run notice")
	}
}

func TestExecuteMission(t *testing.T) {
	t.Run("Primary success", func(t *testing.T) {
		executor := &MockCommandExecutor{
			RunFunc: func(name string, args []string, env []string, stdout, stderr io.Writer) error {
				return nil
			},
		}
		opts := AgentOptions{Executor: executor}
		err := executeMission(opts, "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Primary fails, fallback succeeds", func(t *testing.T) {
		callCount := 0
		executor := &MockCommandExecutor{
			RunFunc: func(name string, args []string, env []string, stdout, stderr io.Writer) error {
				callCount++
				if callCount == 1 {
					return fmt.Errorf("primary failed")
				}
				return nil
			},
		}
		opts := AgentOptions{
			Executor:      executor,
			FallbackModel: "fallback",
		}
		err := executeMission(opts, "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if callCount != 2 {
			t.Errorf("expected 2 calls, got %d", callCount)
		}
	})

	t.Run("Both fail", func(t *testing.T) {
		executor := &MockCommandExecutor{
			RunFunc: func(name string, args []string, env []string, stdout, stderr io.Writer) error {
				return fmt.Errorf("fail")
			},
		}
		opts := AgentOptions{
			Executor:      executor,
			FallbackModel: "fallback",
		}
		err := executeMission(opts, "")
		if err == nil {
			t.Error("expected error when both fail")
		}
	})
}

func TestRealCommandExecutor_RunCommand(t *testing.T) {
	executor := &RealCommandExecutor{}
	var stdout, stderr bytes.Buffer
	err := executor.RunCommand("echo", []string{"hello"}, os.Environ(), &stdout, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout.String(), "hello") {
		t.Errorf("expected stdout to contain hello, got %q", stdout.String())
	}
}

func TestGetEnvOrDefault(t *testing.T) {

	os.Setenv("TEST_VAR", "value")
	defer os.Unsetenv("TEST_VAR")

	if getEnvOrDefault("TEST_VAR", "default") != "value" {
		t.Error("expected value from env")
	}
	if getEnvOrDefault("NON_EXISTENT", "default") != "default" {
		t.Error("expected default value")
	}
}

package cli

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestGeminiCLI_Install(t *testing.T) {
	callCount := 0
	mock := &MockCommandExecutor{
		RunFunc: func(name string, args []string, env []string, stdout, stderr io.Writer) error {
			callCount++
			if name != "sh" {
				return fmt.Errorf("expected sh command")
			}
			return nil
		},
	}

	gemini := &GeminiCLI{}
	err := gemini.Install(mock)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 command call, got %d", callCount)
	}
}

func TestGeminiCLI_Auth(t *testing.T) {
	gemini := &GeminiCLI{}
	err := gemini.Auth(nil, "token")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if gemini.Token != "token" {
		t.Errorf("expected token to be set")
	}
}

func TestGeminiCLI_Run(t *testing.T) {
	mock := &MockCommandExecutor{
		RunFunc: func(name string, args []string, env []string, stdout, stderr io.Writer) error {
			if name != "gemini" {
				return fmt.Errorf("expected gemini command")
			}
			foundToken := false
			for _, e := range env {
				if strings.HasPrefix(e, "GEMINI_API_KEY=") {
					foundToken = true
				}
			}
			if !foundToken {
				return fmt.Errorf("expected GEMINI_API_KEY in env")
			}
			return nil
		},
	}

	gemini := &GeminiCLI{Token: "test-token"}
	err := gemini.Run(mock, "prompt", "model")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

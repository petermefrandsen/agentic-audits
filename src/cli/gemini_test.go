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
			if name != "curl" {
				return fmt.Errorf("expected curl command, got %s", name)
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
			if name != "curl" {
				return fmt.Errorf("expected curl command")
			}
			// Verify arguments contain URL, headers, and fail flag
			foundURL := false
			foundHeader := false
			foundFail := false
			foundShowBlock := false
			for _, arg := range args {
				if strings.Contains(arg, "generateContent") {
					foundURL = true
				}
				if strings.Contains(arg, "x-goog-api-key: test-token") {
					foundHeader = true
				}
				if arg == "-f" {
					foundFail = true
				}
				if arg == "-S" {
					foundShowBlock = true
				}
			}

			if !foundURL {
				return fmt.Errorf("expected URL in args")
			}
			if !foundHeader {
				// We expect token in header now
				return fmt.Errorf("expected api key header in args")
			}
			if !foundFail {
				return fmt.Errorf("expected -f flag in args")
			}
			if !foundShowBlock {
				return fmt.Errorf("expected -S flag in args")
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

package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	// Mock executor and http client
	executor := &MockCommandExecutor{
		RunFunc: func(name string, args []string, env []string, stdout, stderr io.Writer) error {
			return nil
		},
	}
	httpClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"login": "test-user"}`)),
				}, nil
			},
	}

	// Setup tmp config dir
	tmpHome, _ := os.MkdirTemp("", "main-test")
	defer os.RemoveAll(tmpHome)
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", oldHome)

	t.Run("Basic success", func(t *testing.T) {
		err := run([]string{"--mission", "test", "--github-token", "tok", "--skip-setup"}, executor, httpClient)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("With setup success", func(t *testing.T) {
		err := run([]string{"--mission", "test", "--github-token", "tok"}, executor, httpClient)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Invalid flags", func(t *testing.T) {
		err := run([]string{"--invalid"}, executor, httpClient)
		if err == nil {
			t.Error("expected error for invalid flags")
		}
	})
}

func TestOutputEnv(t *testing.T) {
	tmpFile, _ := os.CreateTemp("", "env")
	defer os.Remove(tmpFile.Name())
	os.Setenv("GITHUB_ENV", tmpFile.Name())
	defer os.Unsetenv("GITHUB_ENV")

	outputEnv("TEST_KEY", "test_value")
	outputEnv("MULTI", "line\nvalue")

	data, _ := os.ReadFile(tmpFile.Name())
	if !bytes.Contains(data, []byte("TEST_KEY=test_value")) {
		t.Error("should contain single line env")
	}
	if !bytes.Contains(data, []byte("MULTI<<EOF")) {
		t.Error("should contain multi-line env")
	}
}

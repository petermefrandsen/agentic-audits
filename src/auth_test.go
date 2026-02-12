package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestConfigureGitHubAuth(t *testing.T) {
	// Setup temporary home for .config/gh
	tmpHome, err := os.MkdirTemp("", "home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)
	
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", oldHome)

	t.Run("Empty token", func(t *testing.T) {
		err := configureGitHubAuth(&MockHTTPClient{}, "")
		if err == nil {
			t.Error("expected error for empty token")
		}
	})

	t.Run("Valid token and successful API call", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"login": "test-user"}`)),
				}, nil
			},
		}

		err := configureGitHubAuth(mockClient, "valid-token")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		hostsFile := filepath.Join(tmpHome, ".config", "gh", "hosts.yml")
		data, err := os.ReadFile(hostsFile)
		if err != nil {
			t.Fatalf("failed to read hosts.yml: %v", err)
		}

		if !bytes.Contains(data, []byte(`user: "test-user"`)) {
			t.Error("hosts.yml should contain detected username")
		}
		if !bytes.Contains(data, []byte(`oauth_token: "valid-token"`)) {
			t.Error("hosts.yml should contain token")
		}
	})

	t.Run("API call fails, use default username", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("network error")
			},
		}

		err := configureGitHubAuth(mockClient, "token")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		hostsFile := filepath.Join(tmpHome, ".config", "gh", "hosts.yml")
		data, err := os.ReadFile(hostsFile)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Contains(data, []byte(`user: "headless-agent"`)) {
			t.Error("hosts.yml should contain default username on failure")
		}
	})
}

package main

import (
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

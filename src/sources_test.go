package main

import (
	"os"
	"testing"
)

func TestParseSources(t *testing.T) {
	// Test empty config path
	res, err := parseSources("")
	if err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
	if len(res.MCPServers) != 0 {
		t.Errorf("expected 0 mcp servers, got %d", len(res.MCPServers))
	}

	// Test non-existent file
	res, err = parseSources("non-existent.yml")
	if err != nil {
		t.Errorf("expected no error for non-existent file, got %v", err)
	}

	// Test valid YAML
	content := `
- name: test-mcp
  type: mcp
  package: test-package
  enabled: true
- name: test-web
  type: web
  url: https://example.com
  enabled: true
- name: disabled
  type: mcp
  package: disabled-package
  enabled: false
`
	tmpfile, err := os.CreateTemp("", "sources.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	res, err = parseSources(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to parse valid sources: %v", err)
	}

	if len(res.MCPServers) != 1 {
		t.Errorf("expected 1 mcp server, got %d", len(res.MCPServers))
	}
	if res.MCPServers["test-mcp"].Command != "npx" {
		t.Errorf("expected command npx, got %s", res.MCPServers["test-mcp"].Command)
	}
	if res.MCPPackages[0] != "test-package" {
		t.Errorf("expected package test-package, got %s", res.MCPPackages[0])
	}
	if res.WebSources != "Also consult these documentation sources: https://example.com" {
		t.Errorf("unexpected web sources: %s", res.WebSources)
	}

	// Test invalid YAML
	invalidContent := `
- name: test
  type: [invalid
`
	tmpfile2, err := os.CreateTemp("", "invalid.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile2.Name())
	tmpfile2.Write([]byte(invalidContent))
	tmpfile2.Close()

	_, err = parseSources(tmpfile2.Name())
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestResolveMission(t *testing.T) {
	// Test mission resolution
	m, err := resolveMission("hello", "")
	if err != nil || m != "hello" {
		t.Errorf("failed to resolve mission: %v", err)
	}

	// Test template resolution
	err = os.MkdirAll(".github/templates", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(".github/templates/test.md", []byte("template content"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(".github/templates")

	m, err = resolveMission("", "test")
	if err != nil || m != "template content" {
		t.Errorf("failed to resolve template: %v", err)
	}

	// Test both provided
	_, err = resolveMission("m", "t")
	if err == nil {
		t.Error("expected error when both provided")
	}

	// Test neither provided
	_, err = resolveMission("", "")
	if err == nil {
		t.Error("expected error when neither provided")
	}

	// Test missing template
	_, err = resolveMission("", "missing")
	if err == nil {
		t.Error("expected error for missing template")
	}
}

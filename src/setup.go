package main

import (
	"fmt"
	"os"
	"os/exec"
)

func installGitHubCLI() error {
	if _, err := exec.LookPath("gh"); err == nil {
		fmt.Println("GitHub CLI is already installed.")
		return nil
	}

	fmt.Println("::group::Installing GitHub CLI")
	defer fmt.Println("::endgroup::")

	// This follows the logic from install_gh.sh for Debian-based systems (common in GH Actions)
	commands := [][]string{
		{"curl", "-fsSL", "https://cli.github.com/packages/githubcli-archive-keyring.gpg"},
		{"sudo", "dd", "of=/usr/share/keyrings/githubcli-archive-keyring.gpg"},
		{"sh", "-c", "echo \"deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main\" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null"},
		{"sudo", "apt-get", "update", "-qq"},
		{"sudo", "apt-get", "install", "-y", "-qq", "gh"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run %v: %w", cmdArgs, err)
		}
	}

	return nil
}

func installCopilotExtension() error {
	fmt.Println("::group::Installing gh-copilot extension")
	defer fmt.Println("::endgroup::")

	// Check if already installed
	cmd := exec.Command("gh", "extension", "list")
	output, _ := cmd.CombinedOutput()
	if contains(string(output), "github/gh-copilot") {
		fmt.Println("gh-copilot is already available.")
		return nil
	}

	fmt.Println("Installing github/gh-copilot extension...")
	cmd = exec.Command("gh", "extension", "install", "github/gh-copilot", "--force")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && stringContains(s, substr)
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

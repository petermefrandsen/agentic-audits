package main

import (
	"fmt"
	"os"
)


func installGitHubCLI(executor CommandExecutor) error {
	// We can't easily check for existence via executor without returning something.
	// Let's assume the caller handles basic existence check or the executor does.
	// Actually, let's keep it simple and just run the commands.
	
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
		if err := executor.RunCommand(cmdArgs[0], cmdArgs[1:], os.Environ(), os.Stdout, os.Stderr); err != nil {
			return fmt.Errorf("failed to run %v: %w", cmdArgs, err)
		}
	}


	return nil
}

func installCopilotExtension(executor CommandExecutor) error {
	fmt.Println("::group::Installing gh-copilot extension")
	defer fmt.Println("::endgroup::")

	fmt.Println("Installing github/gh-copilot extension...")
	return executor.RunCommand("gh", []string{"extension", "install", "github/gh-copilot", "--force"}, os.Environ(), os.Stdout, os.Stderr)
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

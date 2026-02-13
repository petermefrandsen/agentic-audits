package cli

import (
	"fmt"
	"os"
)

type CopilotCLI struct {
	Token string
}

func (c *CopilotCLI) Install(executor CommandExecutor) error {
	fmt.Println("::group::Installing GitHub CLI (Copilot Requirement)")

	// Note: We are not checking if it's already installed for simplicity, as per original code.
	commands := [][]string{
		{"curl", "-fsSL", "https://cli.github.com/packages/githubcli-archive-keyring.gpg"},
		{"sudo", "dd", "of=/usr/share/keyrings/githubcli-archive-keyring.gpg"},
		{"sh", "-c", "echo \"deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main\" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null"},
		{"sudo", "apt-get", "update", "-qq"},
		{"sudo", "apt-get", "install", "-y", "-qq", "gh"},
	}

	for _, cmdArgs := range commands {
		if err := executor.RunCommand(cmdArgs[0], cmdArgs[1:], os.Environ(), os.Stdout, os.Stderr); err != nil {
			fmt.Println("::endgroup::")
			return fmt.Errorf("failed to run %v: %w", cmdArgs, err)
		}
	}
	fmt.Println("::endgroup::")

	fmt.Println("::group::Installing gh-copilot extension")
	fmt.Println("Installing github/gh-copilot extension...")
	if err := executor.RunCommand("gh", []string{"extension", "install", "github/gh-copilot", "--force"}, os.Environ(), os.Stdout, os.Stderr); err != nil {
		fmt.Println("::endgroup::")
		return err
	}
	fmt.Println("::endgroup::")

	return nil
}

func (c *CopilotCLI) Auth(executor CommandExecutor, token string) error {
	c.Token = token
	return nil
}

func (c *CopilotCLI) Run(executor CommandExecutor, prompt string, model string) error {
	if c.Token == "" {
		return fmt.Errorf("token not set, call Auth first")
	}

	args := []string{"copilot", "--allow-all-tools", "-p", prompt}
	if model != "" {
		args = append(args, "--model", model)
	}

	env := append(os.Environ(),
		"COPILOT_GITHUB_TOKEN="+c.Token,
		"GITHUB_TOKEN="+c.Token,
	)

	fmt.Printf("Running copilot agent with model: %s\n", model)
	return executor.RunCommand("gh", args, env, os.Stdout, os.Stderr)
}

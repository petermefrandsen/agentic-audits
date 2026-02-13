package cli

import (
	"fmt"
	"os"
)

type GeminiCLI struct {
	Token string
}

func (c *GeminiCLI) Install(executor CommandExecutor) error {
	fmt.Println("::group::Installing Gemini CLI")

	// Commands to install Gemini CLI
	// Using the installation script from geminicli.com
	// curl -sSfL https://raw.githubusercontent.com/google-gemini/gemini-cli/main/install.sh | sh

	cmd := []string{"sh", "-c", "curl -sSfL https://raw.githubusercontent.com/google-gemini/gemini-cli/main/install.sh | sh"}

	if err := executor.RunCommand(cmd[0], cmd[1:], os.Environ(), os.Stdout, os.Stderr); err != nil {
		fmt.Println("::endgroup::")
		return fmt.Errorf("failed to install Gemini CLI: %w", err)
	}

	// Add to PATH if needed? The install script usually handles it or tells you to.
	// For CI/CD, we might need to add it to generic path or call it directly.
	// Assuming it installs to ~/.gemini/bin or similar and we might need to update PATH or call absolute path.
	// Standard install location is often /usr/local/bin or ~/.local/bin

	fmt.Println("::endgroup::")
	return nil
}

func (c *GeminiCLI) Auth(executor CommandExecutor, token string) error {
	c.Token = token
	return nil
}

func (c *GeminiCLI) Run(executor CommandExecutor, prompt string, model string) error {
	if c.Token == "" {
		return fmt.Errorf("token (GEMINI_API_KEY) not set, call Auth first")
	}

	// Command: gemini prompt <prompt>
	// Or headless mode specifics?
	// Docs say: gemini <prompt> or gemini prompt <prompt>
	// We want to use it in a similar way to copilot.

	args := []string{prompt}
	// Model selection?
	if model != "" {
		args = append([]string{"--model", model}, args...)
	}

	// Check if we need --headless flag if it exists, or if it detects it.
	// Based on user request/docs: "headless mode https://geminicli.com/docs/cli/headless/"
	// Usually standard CLI usage is headless if not interactive.

	env := append(os.Environ(),
		"GEMINI_API_KEY="+c.Token,
	)

	fmt.Printf("Running gemini agent with model: %s\n", model)
	// Assuming binary is 'gemini' and is in PATH after install
	return executor.RunCommand("gemini", args, env, os.Stdout, os.Stderr)
}

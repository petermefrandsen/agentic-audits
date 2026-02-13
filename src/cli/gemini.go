package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

type GeminiCLI struct {
	Token string
}

func (c *GeminiCLI) Install(executor CommandExecutor) error {
	fmt.Println("::group::Verifying curl for Gemini")
	// Verify curl exists
	cmd := []string{"curl", "--version"}
	if err := executor.RunCommand(cmd[0], cmd[1:], os.Environ(), os.Stdout, os.Stderr); err != nil {
		fmt.Println("::endgroup::")
		return fmt.Errorf("curl is required for Gemini CLI but was not found: %w", err)
	}
	fmt.Println("::endgroup::")
	return nil
}

func (c *GeminiCLI) Auth(executor CommandExecutor, token string) error {
	c.Token = token
	return nil
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

func (c *GeminiCLI) Run(executor CommandExecutor, prompt string, model string) error {
	if c.Token == "" {
		return fmt.Errorf("token (GEMINI_API_KEY) not set, call Auth first")
	}
	if model == "" {
		model = "gemini-pro"
	}

	// Create request body
	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal gemini request: %w", err)
	}

	// Write to temp file
	tmpFile, err := os.CreateTemp("", "gemini-req-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	fmt.Printf("Running gemini agent with model: %s\n", model)

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", model)
	args := []string{
		"-s",
		"-S", // Show error buffer when -s is used
		"-f", // Fail on server errors
		"-X", "POST",
		url,
		"-H", "Content-Type: application/json",
		"-H", fmt.Sprintf("x-goog-api-key: %s", c.Token),
		"-d", fmt.Sprintf("@%s", tmpFile.Name()),
	}

	// NOTE: If RunCommand logs args, the token is leaked in the header arg!
	// CommandExecutor usually just runs.
	// If the user's executor logs args, we are in trouble.
	// But `RealCommandExecutor` in `agent.go` does not log args.

	// We pass generic environment, but here we are using args for auth.

	return executor.RunCommand("curl", args, os.Environ(), os.Stdout, os.Stderr)
}

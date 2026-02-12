package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type AgentOptions struct {
	FullMission   string
	ContextFiles  string
	Model         string
	FallbackModel string
	DryRun        bool
	GithubToken   string
}

func constructFullPrompt(mission string, options AgentOptions, webSources string) string {
	fullMission := fmt.Sprintf("%s (context files: %s)", mission, options.ContextFiles)
	if webSources != "" {
		fullMission = fmt.Sprintf("%s. %s", fullMission, webSources)
	}

	if !options.DryRun {
		fullMission += fmt.Sprintf(`

### MANDATORY: Pull Request Creation
You MUST create a Pull Request for your changes using the `+"`create_pull_request`"+` tool from the GitHub MCP server. 

PR Specifications:
- **Repository**: %s
- **Base Branch**: %s
- **Branch Name**: %s
- **Title**: %s
- **Body**: %s
- **Labels**: %s
`, 
			os.Getenv("GITHUB_REPOSITORY"),
			getEnvOrDefault("PR_BASE", "main"),
			getEnvOrDefault("PR_BRANCH", fmt.Sprintf("agent/audit-%d", time.Now().Unix())),
			getEnvOrDefault("PR_TITLE", "Use STRICT Conventional Commits format (e.g., refactor(skills): [AI-GENERATED] audit and clarify instructions)."),
			getEnvOrDefault("PR_BODY", `You MUST provide a comprehensive, elite-quality description structured as follows:
### 🔎 Audit Overview
Provide a high-level technical summary of what was audited and the general state of the skills.

### 🛠 Detailed Changes
Provide a per-skill breakdown of specific technical improvements (e.g., Skill X: Removed 40% verbosity, updated paths to match current source tree).

### ⚠️ Manual Review Required
List any specific files where you added <!-- ISSUE --> comments because they require human intervention.`),
			getEnvOrDefault("PR_LABELS", "automated-pr"),
		)
	} else {
		fullMission += `

NOTE: dry_run is set to TRUE. Do NOT create a Pull Request. Just verify the changes and report what you would have done.
`
	}

	return fullMission
}

func runAgent(prompt string, model string, token string) error {
	args := []string{"copilot", "--allow-all-tools", "-p", prompt}
	if model != "" {
		args = append(args, "--model", model)
	}

	cmd := exec.Command("gh", args...)
	cmd.Env = append(os.Environ(), 
		"COPILOT_GITHUB_TOKEN="+token,
		"GITHUB_TOKEN="+token,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running agent with model: %s\n", model)
	return cmd.Run()
}

func executeMission(options AgentOptions, webSources string) error {
	fullPrompt := constructFullPrompt(options.FullMission, options, webSources)

	// Attempt with primary model
	err := runAgent(fullPrompt, options.Model, options.GithubToken)
	if err == nil {
		fmt.Println("Agent mission completed successfully.")
		return nil
	}

	fmt.Printf("::warning::Primary model failed: %v\n", err)

	if options.FallbackModel != "" {
		fmt.Printf("Retrying with fallback model: %s\n", options.FallbackModel)
		err = runAgent(fullPrompt, options.FallbackModel, options.GithubToken)
		if err == nil {
			fmt.Println("Agent mission completed with fallback model.")
			return nil
		}
		return fmt.Errorf("agent mission failed with both primary and fallback models: %w", err)
	}

	return fmt.Errorf("agent mission failed and no fallback model is configured: %w", err)
}

func getEnvOrDefault(name, defaultValue string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return defaultValue
}

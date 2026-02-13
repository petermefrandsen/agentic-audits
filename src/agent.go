package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/petermefrandsen/agentic-audits/src/cli"
)

// CommandExecutor is defined here as well or we can use the one in cli package if we export it.
// To avoid duplication, let's strictly use the one in `cli` package if possible, or alias it.
// But `main` package is special.
// Let's redefine it here or use `cli.CommandExecutor`?
// If we use `cli.CommandExecutor` we need to import `cli`.
// Let's use `cli.CommandExecutor` for consistency.

type RealCommandExecutor struct{}

func (e *RealCommandExecutor) RunCommand(name string, args []string, env []string, stdout, stderr io.Writer) error {
	cmd := exec.Command(name, args...)
	cmd.Env = env
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

type AgentOptions struct {
	FullMission   string
	ContextFiles  string
	Model         string
	FallbackModel string
	DryRun        bool
	GithubToken   string // Still needed for PR creation prompting? Yes.
	Executor      cli.CommandExecutor
	CLI           cli.AICLI
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

func runAgent(cli cli.AICLI, executor cli.CommandExecutor, prompt string, model string) error {
	// The `cli` handles the command construction and env injection
	return cli.Run(executor, prompt, model)
}

func executeMission(options AgentOptions, webSources string) error {
	fullPrompt := constructFullPrompt(options.FullMission, options, webSources)

	// Attempt with primary model
	err := runAgent(options.CLI, options.Executor, fullPrompt, options.Model)
	if err == nil {
		fmt.Println("Agent mission completed successfully.")
		return nil
	}

	fmt.Printf("::warning::Primary model failed: %v\n", err)

	if options.FallbackModel != "" {
		fmt.Printf("Retrying with fallback model: %s\n", options.FallbackModel)
		err = runAgent(options.CLI, options.Executor, fullPrompt, options.FallbackModel)
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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/petermefrandsen/agentic-audits/src/cli"
)

func main() {
	if err := run(os.Args[1:], &RealCommandExecutor{}, http.DefaultClient); err != nil {
		fmt.Printf("::error::%v\n", err)
		os.Exit(1)
	}
}

func run(args []string, executor cli.CommandExecutor, httpClient HTTPClient) error {
	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	mission := fs.String("mission", "", "Agent mission prompt")
	template := fs.String("template", "", "Mission template name")
	sourcesConfig := fs.String("sources-config", ".github/sources.yml", "Path to sources config")
	githubToken := fs.String("github-token", "", "GitHub Token")
	contextFiles := fs.String("context-files", ".", "Context files or globs")
	model := fs.String("model", "", "Primary model")
	fallbackModel := fs.String("fallback-model", "", "Fallback model")
	dryRun := fs.Bool("dry-run", false, "Skip PR creation")
	skipSetup := fs.Bool("skip-setup", false, "Skip CLI and extension installation")
	cliName := fs.String("cli", "copilot", "AI CLI to use (copilot or gemini)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Select CLI
	var aiCLI cli.AICLI
	switch strings.ToLower(*cliName) {
	case "copilot":
		aiCLI = &cli.CopilotCLI{}
	case "gemini":
		aiCLI = &cli.GeminiCLI{}
	default:
		return fmt.Errorf("unsupported CLI: %s", *cliName)
	}

	// 0. Setup (CLI, Auth, Extension)
	if !*skipSetup {
		if err := aiCLI.Install(executor); err != nil {
			fmt.Printf("::warning::Setup failed (%s CLI): %v\n", *cliName, err)
		}

		// Auth
		// We pass githubToken to Auth. Copilot uses it as GITHUB_TOKEN/COPILOT_GITHUB_TOKEN
		// Gemini uses it as GEMINI_API_KEY (if passed via this flag, or we expect env var)
		// For Gemini, strict adherence to implementation plan: pass through env var.
		// But here we might want to allow the flag to serve as the key if provided?
		// The `action.yml` input is `github_token`.
		// If using Gemini, the token input might be the Gemini API Key?
		// Or we expect a separate secret/env var?
		// Implementation plan said: "Gemini CLI authentication for headless mode will primarily use the GEMINI_API_KEY environment variable."
		// So we assume the environment has it.
		// However, `github-token` is required in action.yml.
		// Let's assume `github-token` is for GitHub operations (PR creation etc) regardless of CLI.
		// And Gemini API Key is separate.
		// But wait, `AICLI.Auth(executor, token)` signature.
		// For Copilot, token is GitHub token.
		// For Gemini, user might want to pass key via separate mechanism.
		// But let's check if we can pass a specific token for the CLI.

		// If CLI is Gemini, we probably shouldn't pass the GitHub token to `Auth` if `Auth` expects an API key.
		// UNLESS we repurpose the `github-token` input to be "AI Token".
		// But we still need GitHub token for PR creation!
		// So we likely need `GEMINI_API_KEY` in environment.
		// The `Auth` method might fetch from env if token arg is empty or irrelevant?
		// Or we assume `Auth` takes the "AI Auth Token".

		// For now, let's pass `githubToken` to Copilot.
		// For Gemini, we might pass empty string if it relies on `GEMINI_API_KEY` env var,
		// OR we pass `os.Getenv("GEMINI_API_KEY")`.

		var authToken string
		if *cliName == "copilot" {
			authToken = *githubToken
			// Also configure internal GH auth for PRs?
			// `configureGitHubAuth` in `setup.go` was doing `gh auth login`.
			// We should perhaps keep `configureGitHubAuth` for general GH operations (PRs)
			if err := configureGitHubAuth(httpClient, *githubToken); err != nil {
				return fmt.Errorf("gh auth failed: %w", err)
			}
		} else if *cliName == "gemini" {
			authToken = os.Getenv("GEMINI_API_KEY")
			if authToken == "" {
				fmt.Println("::error::GEMINI_API_KEY is missing or empty in environment!")
			} else {
				fmt.Println("::debug::GEMINI_API_KEY found in environment.")
			}
			// We still need GH auth for PR creation!
			if *githubToken != "" {
				if err := configureGitHubAuth(httpClient, *githubToken); err != nil {
					return fmt.Errorf("gh auth failed: %w", err)
				}
			}
		}

		if err := aiCLI.Auth(executor, authToken); err != nil {
			return fmt.Errorf("ai cli auth failed: %w", err)
		}
	} else {
		// If skipping setup, we still need to set the token for the instance if it needs it for Run
		var authToken string
		if *cliName == "copilot" {
			authToken = *githubToken
		} else if *cliName == "gemini" {
			authToken = os.Getenv("GEMINI_API_KEY")
		}
		aiCLI.Auth(executor, authToken)
	}

	// 1. Resolve Mission
	resolvedMission, err := resolveMission(*mission, *template)
	if err != nil {
		fmt.Printf("::error::%v\n", err)
		os.Exit(1)
	}

	// 2. Configure Sources
	processed, err := parseSources(*sourcesConfig)
	if err != nil {
		fmt.Printf("::error::Error parsing sources: %v\n", err)
	}

	// 3. Write Copilot Config (Only if Copilot?)
	// If using Gemini, do we need this?
	// The `sources` config seems to process MCP servers.
	// Gemini CLI might support MCP?
	// If not, this block might be Copilot specific.
	// For now, let's guard it or keep it if it doesn't hurt.
	// It writes to `~/.config/github-copilot`.
	if *cliName == "copilot" {
		configDir := filepath.Join(os.Getenv("HOME"), ".config", "github-copilot")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("::error::Failed to create config dir: %v\n", err)
			os.Exit(1)
		}

		copilotConfig := CopilotConfig{
			MCPServers: processed.MCPServers,
		}
		configData, _ := json.MarshalIndent(copilotConfig, "", "  ")
		configFile := filepath.Join(configDir, "config.json")
		if err := os.WriteFile(configFile, configData, 0644); err != nil {
			fmt.Printf("::error::Failed to write config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Copilot config written to", configFile)
		fmt.Println(string(configData))
	}

	// 4. Handle Output/Env
	outputEnv("RESOLVED_MISSION", resolvedMission)
	outputEnv("EXTRA_WEB_SOURCES", processed.WebSources)

	// 5. Verify gh
	if _, err := exec.LookPath("gh"); err != nil {
		fmt.Println("::warning::gh CLI not found in path")
	}

	// 6. Execute Mission
	agentOpts := AgentOptions{
		FullMission:   resolvedMission,
		ContextFiles:  *contextFiles,
		Model:         *model,
		FallbackModel: *fallbackModel,
		DryRun:        *dryRun,
		GithubToken:   *githubToken,
		Executor:      executor,
		CLI:           aiCLI,
	}

	if err := executeMission(agentOpts, processed.WebSources); err != nil {
		return fmt.Errorf("mission execution failed: %w", err)
	}

	return nil
}

func outputEnv(name, value string) {
	if value == "" {
		return
	}
	githubEnv := os.Getenv("GITHUB_ENV")
	if githubEnv != "" {
		f, err := os.OpenFile(githubEnv, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("::error::Failed to open GITHUB_ENV: %v\n", err)
			return
		}
		defer f.Close()

		if strings.Contains(value, "\n") {
			fmt.Fprintf(f, "%s<<EOF\n%s\nEOF\n", name, value)
		} else {
			fmt.Fprintf(f, "%s=%s\n", name, value)
		}
	} else {
		// Fallback for local testing
		fmt.Printf("EXPORT %s=%s\n", name, value)
	}
}

func resolveMission(mission, template string) (string, error) {
	if mission != "" && template != "" {
		return "", fmt.Errorf("both 'mission' and 'template' provided")
	}
	if mission == "" && template == "" {
		return "", fmt.Errorf("neither 'mission' nor 'template' provided")
	}

	if template != "" {
		templatePath := filepath.Join(".github", "templates", template+".md")
		data, err := os.ReadFile(templatePath)
		if err != nil {
			return "", fmt.Errorf("template file not found: %s", templatePath)
		}
		return string(data), nil
	}

	return mission, nil
}

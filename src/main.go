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
)


func main() {
	if err := run(os.Args[1:], &RealCommandExecutor{}, http.DefaultClient); err != nil {
		fmt.Printf("::error::%v\n", err)
		os.Exit(1)
	}
}

func run(args []string, executor CommandExecutor, httpClient HTTPClient) error {
	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	mission := fs.String("mission", "", "Agent mission prompt")
	template := fs.String("template", "", "Mission template name")
	sourcesConfig := fs.String("sources-config", ".github/sources.yml", "Path to sources config")
	githubToken := fs.String("github-token", "", "GitHub Token")
	contextFiles := fs.String("context-files", ".", "Context files or globs")
	model := fs.String("model", "", "Primary model")
	fallbackModel := fs.String("fallback-model", "", "Fallback model")
	dryRun := fs.Bool("dry-run", false, "Skip PR creation")
	skipSetup := fs.Bool("skip-setup", false, "Skip GH CLI and extension installation")
	
	if err := fs.Parse(args); err != nil {
		return err
	}

	// 0. Setup (CLI, Auth, Extension)
	if !*skipSetup {
		if err := installGitHubCLI(executor); err != nil {
			fmt.Printf("::warning::Setup failed (GH CLI): %v\n", err)
		}
		if err := configureGitHubAuth(httpClient, *githubToken); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}
		if err := installCopilotExtension(executor); err != nil {
			fmt.Printf("::warning::Setup failed (Copilot extension): %v\n", err)
		}
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
		// Don't exit here, might want to continue without sources? 
		// JS logic returns empty defaults on error.
	}

	// 3. Write Copilot Config
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

	// 4. Handle Output/Env
	outputEnv("RESOLVED_MISSION", resolvedMission)
	outputEnv("EXTRA_WEB_SOURCES", processed.WebSources)

	// 5. Verify gh and optionally run commands
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
	}

	if err := executeMission(agentOpts, processed.WebSources); err != nil {
		return fmt.Errorf("mission execution failed: %w", err)
	}

	// Print summary (mimics ::group:: behavior)
	fmt.Println("Copilot config written to", configFile)
	fmt.Println(string(configData))
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

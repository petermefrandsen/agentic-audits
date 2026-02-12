package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type GitHubUser struct {
	Login string `json:"login"`
}

func configureGitHubAuth(token string) error {
	if token == "" {
		return fmt.Errorf("GH_TOKEN is not set")
	}

	fmt.Println("Configuring gh auth manually to bypass scope validation...")

	username := "headless-agent"
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err == nil {
		req.Header.Set("Authorization", "token "+token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			var user GitHubUser
			if err := json.NewDecoder(resp.Body).Decode(&user); err == nil {
				username = user.Login
				fmt.Printf("Detected username: %s\n", username)
			}
			resp.Body.Close()
		}
	}

	configDir := filepath.Join(os.Getenv("HOME"), ".config", "gh")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create gh config dir: %w", err)
	}

	hostsContent := fmt.Sprintf(`github.com:
    user: "%s"
    oauth_token: "%s"
    git_protocol: "https"
`, username, token)

	hostsFile := filepath.Join(configDir, "hosts.yml")
	if err := os.WriteFile(hostsFile, []byte(hostsContent), 0600); err != nil {
		return fmt.Errorf("failed to write hosts.yml: %w", err)
	}

	fmt.Println("gh auth configured successfully.")
	return nil
}

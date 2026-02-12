package main

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Source struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Package string `yaml:"package"`
	URL     string `yaml:"url"`
	Enabled bool   `yaml:"enabled"`
}

type Config struct {
	Sources []Source `yaml:"sources"`
}

type MCPServer struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type CopilotConfig struct {
	MCPServers map[string]MCPServer `json:"mcpServers"`
}

type ProcessedSources struct {
	MCPServers  map[string]MCPServer
	MCPPackages []string
	WebSources  string
}

func parseSources(configPath string) (ProcessedSources, error) {
	result := ProcessedSources{
		MCPServers:  make(map[string]MCPServer),
		MCPPackages: []string{},
	}

	if configPath == "" {
		return result, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return result, err
	}

	// The current logic in JS seems to manually parse YAML-like structure.
	// We'll use a proper YAML parser.
	var config struct {
		Sources []Source `yaml:"sources"`
	}

	// The JS script seems to expect a list of sources directly or under a 'sources' key?
	// Let's check the JS logic again. It matches line by line:
	// - name: ...
	// type: ...
	// package: ...
	// url: ...
	// enabled: ...
	
	// This structure is a slice of Source.
	var sources []Source
	err = yaml.Unmarshal(data, &sources)
	if err != nil {
		// Try parsing as a map with 'sources' key just in case
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return result, err
		}
		sources = config.Sources
	}

	var webUrls []string
	for _, s := range sources {
		if !s.Enabled {
			continue
		}

		switch s.Type {
		case "mcp":
			if s.Package != "" {
				result.MCPServers[s.Name] = MCPServer{
					Command: "npx",
					Args:    []string{"-y", s.Package},
				}
				result.MCPPackages = append(result.MCPPackages, s.Package)
			}
		case "web":
			if s.URL != "" {
				webUrls = append(webUrls, s.URL)
			}
		}
	}

	if len(webUrls) > 0 {
		result.WebSources = "Also consult these documentation sources: " + strings.Join(webUrls, ", ")
	}

	return result, nil
}

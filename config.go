package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads the configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	path := resolveConfigPath(configPath)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set default separator if not specified
	if config.Separator == "" {
		config.Separator = " | "
	}

	// Validate actions
	if err := validateActions(config.Actions); err != nil {
		return nil, err
	}

	return &config, nil
}

// resolveConfigPath resolves the configuration file path
func resolveConfigPath(configPath string) string {
	// If explicit path is provided, use it
	if configPath != "" {
		return configPath
	}

	// Try XDG_CONFIG_HOME
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		path := filepath.Join(xdgConfigHome, "ccstatusline", "config.yaml")
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Default to ~/.config/ccstatusline/config.yaml
	homeDir, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(homeDir, ".config", "ccstatusline", "config.yaml")
	}

	return "config.yaml"
}

// validateActions validates the action configuration
func validateActions(actions []Action) error {
	names := make(map[string]bool)

	for i, action := range actions {
		// Check name is required
		if action.Name == "" {
			return fmt.Errorf("action at index %d: name is required", i)
		}

		// Check for duplicate names
		if names[action.Name] {
			return fmt.Errorf("duplicate action name: %s", action.Name)
		}
		names[action.Name] = true

		// Check command is required
		if action.Command == "" {
			return fmt.Errorf("action %s: command is required", action.Name)
		}
	}

	return nil
}

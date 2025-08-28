package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `actions:
  - name: test1
    command: "echo 'Test 1'"
    color: cyan
  - name: test2
    command: "echo hello"
    color: green
  - name: test3
    command: "echo 'User:' $(whoami)"
separator: " | "`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test loading the config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Verify the loaded config
	if len(config.Actions) != 3 {
		t.Errorf("Expected 3 actions, got %d", len(config.Actions))
	}

	// Check first action (text only)
	if config.Actions[0].Name != "test1" {
		t.Errorf("First action name = %q, want %q", config.Actions[0].Name, "test1")
	}
	if config.Actions[0].Command != "echo 'Test 1'" {
		t.Errorf("First action command = %q, want %q", config.Actions[0].Command, "echo 'Test 1'")
	}
	if config.Actions[0].Color != "cyan" {
		t.Errorf("First action color = %q, want %q", config.Actions[0].Color, "cyan")
	}

	// Check second action (command only)
	if config.Actions[1].Command != "echo hello" {
		t.Errorf("Second action command = %q, want %q", config.Actions[1].Command, "echo hello")
	}
	if config.Actions[1].Color != "green" {
		t.Errorf("Second action color = %q, want %q", config.Actions[1].Color, "green")
	}

	// Check third action (command with template)
	if config.Actions[2].Command != "echo 'User:' $(whoami)" {
		t.Errorf("Third action command = %q, want %q", config.Actions[2].Command, "echo 'User:' $(whoami)")
	}

	if config.Separator != " | " {
		t.Errorf("Separator = %q, want %q", config.Separator, " | ")
	}
}

func TestLoadConfigDefaultSeparator(t *testing.T) {
	// Create a config without separator
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `actions:
  - name: test
    command: "echo 'Test'"`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Should default to " | "
	if config.Separator != " | " {
		t.Errorf("Default separator = %q, want %q", config.Separator, " | ")
	}
}

func TestResolveConfigPath(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		setup      func() string
		cleanup    func()
		expected   func(string) string
	}{
		{
			name:       "explicit path",
			configPath: "/explicit/path/config.yaml",
			setup:      func() string { return "" },
			cleanup:    func() {},
			expected:   func(s string) string { return "/explicit/path/config.yaml" },
		},
		{
			name:       "XDG_CONFIG_HOME",
			configPath: "",
			setup: func() string {
				xdgHome := "/tmp/xdg-config"
				os.Setenv("XDG_CONFIG_HOME", xdgHome)
				// Create the file so stat succeeds
				os.MkdirAll(filepath.Join(xdgHome, "ccstatusline"), 0755)
				os.WriteFile(filepath.Join(xdgHome, "ccstatusline", "config.yaml"), []byte("test"), 0644)
				return xdgHome
			},
			cleanup: func() {
				os.Unsetenv("XDG_CONFIG_HOME")
			},
			expected: func(xdgHome string) string {
				return filepath.Join(xdgHome, "ccstatusline", "config.yaml")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupResult := tt.setup()
			defer tt.cleanup()

			result := resolveConfigPath(tt.configPath)
			expected := tt.expected(setupResult)

			if result != expected {
				t.Errorf("resolveConfigPath(%q) = %q, want %q", tt.configPath, result, expected)
			}
		})
	}
}

func TestActionTypeField(t *testing.T) {
	t.Run("parse builtin type action", func(t *testing.T) {
		yamlContent := []byte(`
actions:
  - name: context_percent
    type: builtin
    function: context_percent
    cache_ttl: 3
`)
		config, err := parseConfig(yamlContent)
		if err != nil {
			t.Fatalf("failed to parse config: %v", err)
		}

		if len(config.Actions) != 1 {
			t.Fatalf("expected 1 action, got %d", len(config.Actions))
		}

		action := config.Actions[0]
		if action.Type != "builtin" {
			t.Errorf("expected type 'builtin', got '%s'", action.Type)
		}
		if action.Function != "context_percent" {
			t.Errorf("expected function 'context_percent', got '%s'", action.Function)
		}
	})

	t.Run("parse command type action", func(t *testing.T) {
		yamlContent := []byte(`
actions:
  - name: git_branch
    type: command
    command: git branch --show-current
`)
		config, err := parseConfig(yamlContent)
		if err != nil {
			t.Fatalf("failed to parse config: %v", err)
		}

		action := config.Actions[0]
		if action.Type != "command" {
			t.Errorf("expected type 'command', got '%s'", action.Type)
		}
		if action.Command != "git branch --show-current" {
			t.Errorf("expected command 'git branch --show-current', got '%s'", action.Command)
		}
	})

	t.Run("default type is command", func(t *testing.T) {
		yamlContent := []byte(`
actions:
  - name: git_branch
    command: git branch --show-current
`)
		config, err := parseConfig(yamlContent)
		if err != nil {
			t.Fatalf("failed to parse config: %v", err)
		}

		action := config.Actions[0]
		// Type が指定されていない場合は空文字列になる
		if action.Type != "" {
			t.Errorf("expected type '', got '%s'", action.Type)
		}
	})
}

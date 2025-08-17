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

	// Should default to single space
	if config.Separator != " " {
		t.Errorf("Default separator = %q, want %q", config.Separator, " ")
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

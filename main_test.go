package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainIntegrationSimple(t *testing.T) {
	// Create test config with new structure
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `actions:
  - name: model
    command: "echo '{.model.display_name}'"
    color: cyan
  - name: directory
    command: "echo '{.cwd | split(\"/\") | .[-1]}'"
    color: blue
  - name: session
    command: "echo '[{.session_id | .[0:8]}]'"
    color: gray
separator: " | "`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Prepare test input with correct field names
	inputData := map[string]interface{}{
		"hook_event_name": "Status",
		"model": map[string]interface{}{
			"id":           "claude-3-5",
			"display_name": "Claude 3.5 Sonnet",
		},
		"cwd":             "/home/user/projects/test",
		"session_id":      "abc123def456ghi789",
		"transcript_path": "/tmp/transcript.json",
		"workspace": map[string]interface{}{
			"current_dir": "/home/user/projects/test",
			"project_dir": "/home/user/projects",
		},
		"version": "0.1.0",
		"output_style": map[string]interface{}{
			"name": "default",
		},
	}

	// Run the processor
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	processor := NewProcessor(inputData)
	output, err := processor.Process(config)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	// Check output contains expected parts
	expectedParts := []string{
		"Claude 3.5 Sonnet",
		"test",
		"[abc123de]",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Errorf("Output does not contain expected part %q", part)
		}
	}

	// Check color codes are present
	if !strings.Contains(output, "\033[") {
		t.Error("Output does not contain ANSI color codes")
	}
}

func TestMainWithCommandAction(t *testing.T) {
	// Create test config with command action
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `actions:
  - name: echo_test
    command: "echo 'Hello World'"
  - name: echo_with_template
    command: "echo 'Output: test'"
separator: " | "`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	inputData := map[string]interface{}{
		"hook_event_name": "Status",
		"session_id":      "test123",
		"cwd":             "/test",
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	processor := NewProcessor(inputData)
	output, err := processor.Process(config)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	expected := "Hello World | Output: test"
	if output != expected {
		t.Errorf("Output = %q, want %q", output, expected)
	}
}

// TestErrorHandling tests error handling without calling main()
func TestErrorHandling(t *testing.T) {
	t.Run("invalid JSON input", func(t *testing.T) {
		invalidJSON := "not valid json"
		var inputData map[string]interface{}
		err := json.Unmarshal([]byte(invalidJSON), &inputData)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})

	t.Run("missing config file", func(t *testing.T) {
		_, err := LoadConfig("/nonexistent/config.yaml")
		if err == nil {
			t.Error("Expected error for missing config file")
		}
	})
}

func TestProcessTemplateEdgeCasesSimple(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			name:     "empty template",
			template: "",
			data:     map[string]interface{}{},
			expected: "",
		},
		{
			name:     "no placeholders",
			template: "Static text",
			data:     map[string]interface{}{"key": "value"},
			expected: "Static text",
		},
		{
			name:     "nested missing field",
			template: "{.a.b.c.d}",
			data:     map[string]interface{}{"a": map[string]interface{}{}},
			expected: "",
		},
		{
			name:     "invalid jq syntax",
			template: "{.field | invalid syntax}",
			data:     map[string]interface{}{"field": "value"},
			expected: "[ERROR:",
		},
		{
			name:     "correct field access",
			template: "{.model.display_name} - {.cwd}",
			data: map[string]interface{}{
				"model": map[string]interface{}{
					"display_name": "Opus",
				},
				"cwd": "/home/test",
			},
			expected: "Opus - /home/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processTemplate(tt.template, tt.data)
			if tt.expected == "[ERROR:" {
				if !strings.HasPrefix(result, tt.expected) {
					t.Errorf("processTemplate() = %q, want prefix %q", result, tt.expected)
				}
			} else if result != tt.expected {
				t.Errorf("processTemplate() = %q, want %q", result, tt.expected)
			}
		})
	}
}

package main

import (
	"strings"
	"testing"
)

func TestProcessorSimple(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		inputData map[string]interface{}
		expected  string
		contains  []string
	}{
		{
			name: "simple text output",
			config: &Config{
				Actions: []Action{
					{
						Name:    "test",
						Command: "echo 'Hello World'",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "Hello World",
		},
		{
			name: "text with template",
			config: &Config{
				Actions: []Action{
					{
						Name:    "model",
						Command: "echo 'Model: {.model.display_name}'",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{
				"model": map[string]interface{}{
					"display_name": "Claude 3.5",
					"id":           "claude-3-5",
				},
			},
			expected: "Model: Claude 3.5",
		},
		{
			name: "text with color",
			config: &Config{
				Actions: []Action{
					{
						Name:    "colored",
						Command: "echo 'Status'",
						Color:   "green",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "\033[32mStatus\033[0m",
		},
		{
			name: "multiple actions with separator",
			config: &Config{
				Actions: []Action{
					{
						Name:    "first",
						Command: "echo 'First'",
					},
					{
						Name:    "second",
						Command: "echo 'Second'",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "First | Second",
		},
		{
			name: "command execution",
			config: &Config{
				Actions: []Action{
					{
						Name:    "echo_test",
						Command: "echo 'test-output'",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "test-output",
		},
		{
			name: "command with prefix",
			config: &Config{
				Actions: []Action{
					{
						Name:    "echo_with_prefix",
						Command: "echo 'Result: hello'",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "Result: hello",
		},
		{
			name: "real world example",
			config: &Config{
				Actions: []Action{
					{
						Name:    "model",
						Command: "echo '{.model.display_name}'",
						Color:   "cyan",
					},
					{
						Name:    "git",
						Command: "echo '(main)'", // Mock git command
						Color:   "green",
					},
					{
						Name:    "dir",
						Command: "echo '{.cwd | split(\"/\") | .[-1]}'",
						Color:   "blue",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{
				"model": map[string]interface{}{
					"display_name": "Opus",
				},
				"cwd": "/Users/test/projects/myapp",
			},
			contains: []string{
				"Opus",
				"main",
				"myapp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor(tt.inputData)
			result, err := processor.Process(tt.config)
			if err != nil {
				t.Fatalf("Process() error = %v", err)
			}

			if tt.expected != "" && result != tt.expected {
				t.Errorf("Process() = %q, want %q", result, tt.expected)
			}

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Process() result does not contain %q", substr)
				}
			}
		})
	}
}

func TestProcessorWithCorrectFields(t *testing.T) {
	// Test with actual Claude Code field names
	config := &Config{
		Actions: []Action{
			{
				Name:    "model",
				Command: "echo '{.model.display_name}'",
			},
			{
				Name:    "session",
				Command: "echo '{.session_id | .[0:8]}'",
			},
			{
				Name:    "cwd",
				Command: "echo '{.cwd | split(\"/\") | .[-1]}'",
			},
		},
		Separator: " - ",
	}

	inputData := map[string]interface{}{
		"hook_event_name": "Status",
		"session_id":      "abc123def456789",
		"transcript_path": "/tmp/transcript.json",
		"cwd":             "/Users/test/work/project",
		"model": map[string]interface{}{
			"id":           "claude-opus-4-1",
			"display_name": "Opus",
		},
		"workspace": map[string]interface{}{
			"current_dir": "/Users/test/work/project",
			"project_dir": "/Users/test/work",
		},
		"version": "0.1.0",
		"output_style": map[string]interface{}{
			"name": "default",
		},
	}

	processor := NewProcessor(inputData)
	result, err := processor.Process(config)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	expected := "Opus - abc123de - project"
	if result != expected {
		t.Errorf("Process() = %q, want %q", result, expected)
	}
}

func TestProcessorWithComplexCommand(t *testing.T) {
	// Test complex command pipeline like cchook
	config := &Config{
		Actions: []Action{
			{
				Name:    "transcript_path",
				Command: "cat | jq -r '.transcript_path'",
			},
			{
				Name:    "session_from_stdin",
				Command: "cat | jq -r '.session_id' | cut -c1-8",
			},
		},
		Separator: " - ",
	}

	inputData := map[string]interface{}{
		"session_id":      "abc123def456789",
		"transcript_path": "/tmp/transcript.json",
	}

	processor := NewProcessor(inputData)
	result, err := processor.Process(config)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	expected := "/tmp/transcript.json - abc123de"
	if result != expected {
		t.Errorf("Process() = %q, want %q", result, expected)
	}
}

func TestProcessorWithPrefix(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		inputData map[string]interface{}
		expected  string
	}{
		{
			name: "command with prefix",
			config: &Config{
				Actions: []Action{
					{
						Name:    "session",
						Command: "echo 'abc123'",
						Prefix:  "Session:",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "Session:abc123",
		},
		{
			name: "command with prefix and color",
			config: &Config{
				Actions: []Action{
					{
						Name:    "session",
						Command: "echo 'abc123'",
						Prefix:  "Session:",
						Color:   "green",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "\033[32mSession:abc123\033[0m",
		},
		{
			name: "empty command result should not show prefix",
			config: &Config{
				Actions: []Action{
					{
						Name:    "empty",
						Command: "printf ''",
						Prefix:  "Prefix:",
						Color:   "red",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "",
		},
		{
			name: "failed command should not show prefix",
			config: &Config{
				Actions: []Action{
					{
						Name:    "failed",
						Command: "false",
						Prefix:  "Prefix:",
						Color:   "red",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "",
		},
		{
			name: "multiple actions with prefix",
			config: &Config{
				Actions: []Action{
					{
						Name:    "model",
						Command: "echo 'Claude 3.5'",
						Prefix:  "Model:",
						Color:   "cyan",
					},
					{
						Name:    "session",
						Command: "echo 'xyz789'",
						Prefix:  "ID:",
						Color:   "green",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "\033[36mModel:Claude 3.5\033[0m | \033[32mID:xyz789\033[0m",
		},
		{
			name: "prefix with space",
			config: &Config{
				Actions: []Action{
					{
						Name:    "session",
						Command: "echo 'abc123'",
						Prefix:  "Session: ",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "Session: abc123",
		},
		{
			name: "prefix with template in command",
			config: &Config{
				Actions: []Action{
					{
						Name:    "session",
						Command: "echo '{.session_id}'",
						Prefix:  "ID:",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{
				"session_id": "test-123",
			},
			expected: "ID:test-123",
		},
		{
			name: "mix of actions with and without prefix",
			config: &Config{
				Actions: []Action{
					{
						Name:    "first",
						Command: "echo 'no-prefix'",
					},
					{
						Name:    "second",
						Command: "echo 'with-prefix'",
						Prefix:  "Prefixed:",
					},
					{
						Name:    "third",
						Command: "echo 'also-no-prefix'",
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "no-prefix | Prefixed:with-prefix | also-no-prefix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor(tt.inputData)
			result, err := processor.Process(tt.config)
			if err != nil {
				t.Fatalf("Process() error = %v", err)
			}

			if result != tt.expected {
				t.Errorf("Process() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestProcessorWithCacheSeparationByDirectory(t *testing.T) {
	// 異なるディレクトリからの入力データを作成
	inputData1 := map[string]interface{}{
		"session_id": "test123",
		"cwd":        "/Users/user/work/project1",
		"model": map[string]interface{}{
			"display_name": "Claude 3.5",
		},
	}
	
	inputData2 := map[string]interface{}{
		"session_id": "test456",
		"cwd":        "/Users/user/work/project2",
		"model": map[string]interface{}{
			"display_name": "Claude 3.5",
		},
	}
	
	config := &Config{
		Actions: []Action{
			{
				Name:     "test_action",
				Command:  "echo project1_data",
				CacheTTL: 60,
			},
		},
		Separator: " | ",
	}
	
	// processor1 でキャッシュを設定
	processor1 := NewProcessor(inputData1)
	processor1.cache = NewCache(t.TempDir())
	
	output1, err := processor1.Process(config)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if output1 != "project1_data" {
		t.Errorf("First execution output = %v, want project1_data", output1)
	}
	
	// processor2 は同じキャッシュディレクトリを使うが、異なるcwdなので異なるキャッシュキーになる
	processor2 := NewProcessor(inputData2)
	processor2.cache = processor1.cache // 同じキャッシュインスタンスを共有
	
	// processor2 のコマンドを異なる出力に変更（キャッシュが分離されていることを確認）
	config2 := &Config{
		Actions: []Action{
			{
				Name:     "test_action", // 同じアクション名
				Command:  "echo project2_data",
				CacheTTL: 60,
			},
		},
		Separator: " | ",
	}
	
	output2, err := processor2.Process(config2)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if output2 != "project2_data" {
		t.Errorf("Second execution output = %v, want project2_data", output2)
	}
	
	// processor1 のキャッシュがまだ有効であることを確認
	// コマンドを変更してもキャッシュから読まれるはず
	config3 := &Config{
		Actions: []Action{
			{
				Name:     "test_action",
				Command:  "echo should_not_execute",
				CacheTTL: 60,
			},
		},
		Separator: " | ",
	}
	
	output3, err := processor1.Process(config3)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if output3 != "project1_data" {
		t.Errorf("Cached output for project1 = %v, want project1_data", output3)
	}
}

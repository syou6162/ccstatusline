package builtin

import (
	"os"
	"testing"
)

func TestContextUsedFunction(t *testing.T) {
	fn := ContextUsedFunction{}

	// Name()メソッドのテスト
	if fn.Name() != "context_used" {
		t.Errorf("expected name 'context_used', got '%s'", fn.Name())
	}

	t.Run("valid JSONL transcript with token counts", func(t *testing.T) {
		// テスト用のJSONLトランスクリプトファイルを作成
		transcriptContent := `{"type": "user", "message": {"role": "user", "content": [{"type": "text", "text": "Hello"}]}}
{"type": "assistant", "message": {"role": "assistant", "usage": {"input_tokens": 30000, "cache_creation_input_tokens": 20000, "cache_read_input_tokens": 12000, "output_tokens": 5000}}}
`

		// テンポラリファイルにトランスクリプトを書き込み
		tmpFile := t.TempDir() + "/transcript.jsonl"
		if err := os.WriteFile(tmpFile, []byte(transcriptContent), 0644); err != nil {
			t.Fatal(err)
		}

		input := map[string]interface{}{
			"transcript_path": tmpFile,
		}

		result, err := fn.Execute(input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		// 30000 + 20000 + 12000 = 62000トークン = 62k
		if result != "62k" {
			t.Errorf("expected '62k', got '%s'", result)
		}
	})
}

func TestContextPercentFunction(t *testing.T) {
	fn := ContextPercentFunction{}

	// Name()メソッドのテスト
	if fn.Name() != "context_percent" {
		t.Errorf("expected name 'context_percent', got '%s'", fn.Name())
	}

	t.Run("non-existent transcript file", func(t *testing.T) {
		input := map[string]interface{}{
			"transcript_path": "/tmp/test_transcript.json",
		}

		result, err := fn.Execute(input)
		if err == nil {
			t.Error("expected error for non-existent transcript file, got nil")
		}
		if result != "" {
			t.Errorf("expected empty result for error case, got '%s'", result)
		}
	})

	t.Run("valid JSONL transcript with token counts", func(t *testing.T) {
		// テスト用のJSONLトランスクリプトファイルを作成
		transcriptContent := `{"type": "user", "message": {"role": "user", "content": [{"type": "text", "text": "Hello"}]}}
{"type": "assistant", "message": {"role": "assistant", "usage": {"input_tokens": 10000, "cache_creation_input_tokens": 15000, "cache_read_input_tokens": 5000, "output_tokens": 2000}}}
`

		// テンポラリファイルにトランスクリプトを書き込み
		tmpFile := t.TempDir() + "/transcript.jsonl"
		if err := os.WriteFile(tmpFile, []byte(transcriptContent), 0644); err != nil {
			t.Fatal(err)
		}

		input := map[string]interface{}{
			"transcript_path": tmpFile,
			"model": map[string]interface{}{
				"display_name": "Claude 3.5 Sonnet",
			},
		}

		result, err := fn.Execute(input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		// (10000 + 15000 + 5000) / 200000 * 100 = 15%
		if result != "15%" {
			t.Errorf("expected '15%%', got '%s'", result)
		}
	})
}

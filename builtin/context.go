package builtin

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ContextPercentFunction はコンテキスト使用率を％で返す組み込み関数
type ContextPercentFunction struct{}

// Name は関数名を返す
func (c ContextPercentFunction) Name() string {
	return "context_percent"
}

// Execute はコンテキスト使用率を計算して％付きで返す
func (c ContextPercentFunction) Execute(input map[string]interface{}) (string, error) {
	totalUsed, err := calculateTotalTokens(input)
	if err != nil {
		return "", err
	}

	// モデルの最大トークン数（とりあえず200Kで固定）
	maxTokens := 200000

	// パーセンテージを計算
	percentage := (totalUsed * 100) / maxTokens

	return fmt.Sprintf("%d%%", percentage), nil
}

// ContextUsedFunction はコンテキスト使用量を返す組み込み関数
type ContextUsedFunction struct{}

// Name は関数名を返す
func (c ContextUsedFunction) Name() string {
	return "context_used"
}

// Execute はコンテキスト使用量を計算して返す（例: "62k"）
func (c ContextUsedFunction) Execute(input map[string]interface{}) (string, error) {
	totalUsed, err := calculateTotalTokens(input)
	if err != nil {
		return "", err
	}

	// k表記に変換
	if totalUsed >= 1000 {
		return fmt.Sprintf("%dk", totalUsed/1000), nil
	}
	return fmt.Sprintf("%d", totalUsed), nil
}

// calculateTotalTokens はトランスクリプトファイルからトークン総数を計算する共通関数
// ccusageと同じ方式：最後のassistantメッセージのみを使用
func calculateTotalTokens(input map[string]interface{}) (int, error) {
	// transcript_path を取得
	transcriptPath, ok := input["transcript_path"].(string)
	if !ok || transcriptPath == "" {
		return 0, fmt.Errorf("transcript_path not found in input")
	}

	// ファイル全体を読み込む
	data, err := os.ReadFile(transcriptPath)
	if err != nil {
		return 0, fmt.Errorf("transcript file not found")
	}

	// 行に分割して逆順でイテレート（最新のメッセージから）
	lines := strings.Split(string(data), "\n")

	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// JSON行をパース
		var record struct {
			Type    string `json:"type"`
			Message struct {
				Usage struct {
					InputTokens              int `json:"input_tokens"`
					CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
					CacheReadInputTokens     int `json:"cache_read_input_tokens"`
				} `json:"usage"`
			} `json:"message"`
		}

		if err := json.Unmarshal([]byte(line), &record); err != nil {
			// パースできない行はスキップ
			continue
		}

		// 最初に見つかったassistantメッセージのトークン数を返す
		// ccusageと同じ方式：input + cache_creation + cache_read
		// （outputは含めない - コンテキストには含まれないため）
		if record.Type == "assistant" && record.Message.Usage.InputTokens > 0 {
			totalUsed := record.Message.Usage.InputTokens
			totalUsed += record.Message.Usage.CacheCreationInputTokens
			totalUsed += record.Message.Usage.CacheReadInputTokens
			return totalUsed, nil
		}
	}

	// assistantメッセージが見つからない場合は0を返す
	return 0, nil
}

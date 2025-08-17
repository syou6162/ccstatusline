package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestCache_Get_NotExists(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewCache(tempDir)

	result, ok := cache.Get("test_action")
	if ok {
		t.Error("Expected cache miss, got hit")
	}
	if result != "" {
		t.Errorf("Expected empty result, got: %s", result)
	}
}

func TestCache_Get_Expired(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewCache(tempDir)

	// 期限切れのキャッシュを手動で作成
	cacheFile := filepath.Join(tempDir, "test_action.json")
	expiredData := `{"result":"expired data","expires_at":1}`
	err := os.WriteFile(cacheFile, []byte(expiredData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	result, ok := cache.Get("test_action")
	if ok {
		t.Error("Expected cache miss for expired entry, got hit")
	}
	if result != "" {
		t.Errorf("Expected empty result for expired entry, got: %s", result)
	}
}

func TestCache_Get_Valid(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewCache(tempDir)

	// 有効なキャッシュを手動で作成
	futureTime := time.Now().Add(time.Hour).Unix()
	cacheFile := filepath.Join(tempDir, "test_action.json")
	validData := fmt.Sprintf(`{"result":"valid data","expires_at":%d}`, futureTime)
	err := os.WriteFile(cacheFile, []byte(validData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	result, ok := cache.Get("test_action")
	if !ok {
		t.Error("Expected cache hit, got miss")
	}
	if result != "valid data" {
		t.Errorf("Expected 'valid data', got: %s", result)
	}
}

func TestCache_Set(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewCache(tempDir)

	err := cache.Set("test_action", "test result", 60)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// キャッシュファイルが作成されたか確認
	cacheFile := filepath.Join(tempDir, "test_action.json")
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		t.Error("Cache file was not created")
	}

	// すぐに取得して確認
	result, ok := cache.Get("test_action")
	if !ok {
		t.Error("Expected cache hit after set, got miss")
	}
	if result != "test result" {
		t.Errorf("Expected 'test result', got: %s", result)
	}
}

func TestCache_CleanExpired(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewCache(tempDir)

	// 期限切れファイルを作成
	expiredFile := filepath.Join(tempDir, "expired.json")
	expiredData := `{"result":"expired","expires_at":1}`
	err := os.WriteFile(expiredFile, []byte(expiredData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 有効なファイルを作成
	futureTime := time.Now().Add(time.Hour).Unix()
	validFile := filepath.Join(tempDir, "valid.json")
	validData := fmt.Sprintf(`{"result":"valid","expires_at":%d}`, futureTime)
	err = os.WriteFile(validFile, []byte(validData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// クリーンアップ実行
	err = cache.CleanExpired()
	if err != nil {
		t.Fatalf("Failed to clean expired cache: %v", err)
	}

	// 期限切れファイルが削除されたか確認
	if _, err := os.Stat(expiredFile); !os.IsNotExist(err) {
		t.Error("Expired cache file was not deleted")
	}

	// 有効なファイルが残っているか確認
	if _, err := os.Stat(validFile); os.IsNotExist(err) {
		t.Error("Valid cache file was deleted")
	}
}

func TestCache_GenerateCacheKey(t *testing.T) {
	tests := []struct {
		name        string
		cwd         string
		actionName  string
		wantPattern string // 期待するパターン（正規表現）
	}{
		{
			name:        "simple project path",
			cwd:         "/Users/yasuhisa.yoshida/work/ccstatusline",
			actionName:  "github_pr",
			wantPattern: `^ccstatusline_[a-f0-9]{4}_github_pr$`,
		},
		{
			name:        "different parent same project name",
			cwd:         "/home/user/projects/ccstatusline",
			actionName:  "github_pr",
			wantPattern: `^ccstatusline_[a-f0-9]{4}_github_pr$`,
		},
		{
			name:        "root directory project",
			cwd:         "/ccstatusline",
			actionName:  "git_branch",
			wantPattern: `^ccstatusline_[a-f0-9]{4}_git_branch$`,
		},
	}

	cache := NewCache(t.TempDir())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cache.GenerateCacheKey(tt.cwd, tt.actionName)

			// パターンマッチをチェック
			matched, err := regexp.MatchString(tt.wantPattern, got)
			if err != nil {
				t.Fatalf("Invalid pattern: %v", err)
			}
			if !matched {
				t.Errorf("GenerateCacheKey() = %v, want pattern %v", got, tt.wantPattern)
			}
		})
	}
}

func TestCache_SameProjectDifferentParent(t *testing.T) {
	cache := NewCache(t.TempDir())

	// 同じプロジェクト名だが親ディレクトリが異なる場合
	key1 := cache.GenerateCacheKey("/Users/user1/work/myproject", "action1")
	key2 := cache.GenerateCacheKey("/Users/user2/work/myproject", "action1")

	if key1 == key2 {
		t.Errorf("Same project name with different parents should generate different keys: %v == %v", key1, key2)
	}

	// プロジェクト名部分は同じであることを確認
	project1 := strings.Split(key1, "_")[0]
	project2 := strings.Split(key2, "_")[0]

	if project1 != project2 {
		t.Errorf("Project name part should be the same: %v != %v", project1, project2)
	}
}

func TestCache_GetSet_WithCwd(t *testing.T) {
	cache := NewCache(t.TempDir())

	// 異なるディレクトリで同じアクション名のキャッシュ
	cwd1 := "/Users/user/work/project1"
	cwd2 := "/Users/user/work/project2"
	actionName := "github_pr"

	// project1 のキャッシュを設定
	err := cache.SetWithCwd(cwd1, actionName, "PR-123", 60)
	if err != nil {
		t.Fatalf("SetWithCwd() error = %v", err)
	}

	// project2 のキャッシュを設定
	err = cache.SetWithCwd(cwd2, actionName, "PR-456", 60)
	if err != nil {
		t.Fatalf("SetWithCwd() error = %v", err)
	}

	// project1 のキャッシュを取得
	got1, ok := cache.GetWithCwd(cwd1, actionName)
	if !ok {
		t.Fatal("GetWithCwd() for project1 should return true")
	}
	if got1 != "PR-123" {
		t.Errorf("GetWithCwd() for project1 = %v, want PR-123", got1)
	}

	// project2 のキャッシュを取得
	got2, ok := cache.GetWithCwd(cwd2, actionName)
	if !ok {
		t.Fatal("GetWithCwd() for project2 should return true")
	}
	if got2 != "PR-456" {
		t.Errorf("GetWithCwd() for project2 = %v, want PR-456", got2)
	}
}

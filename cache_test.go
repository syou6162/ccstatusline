package main

import (
	"fmt"
	"os"
	"path/filepath"
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

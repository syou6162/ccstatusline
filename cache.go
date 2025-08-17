package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Cache struct {
	dir string
}

type cacheEntry struct {
	Result    string `json:"result"`
	ExpiresAt int64  `json:"expires_at"`
}

func NewCache(dir string) *Cache {
	return &Cache{dir: dir}
}

func NewDefaultCache() *Cache {
	// XDG Base Directory仕様に従う
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// フォールバック
			cacheDir = filepath.Join(os.TempDir(), "ccstatusline-cache")
		} else {
			cacheDir = filepath.Join(homeDir, ".cache")
		}
	}
	cacheDir = filepath.Join(cacheDir, "ccstatusline")
	return &Cache{dir: cacheDir}
}

func (c *Cache) Get(name string) (string, bool) {
	filePath := filepath.Join(c.dir, name+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", false
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", false
	}

	if time.Now().Unix() > entry.ExpiresAt {
		return "", false
	}

	return entry.Result, true
}

func (c *Cache) Set(name string, result string, ttl int) error {
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return err
	}

	entry := cacheEntry{
		Result:    result,
		ExpiresAt: time.Now().Unix() + int64(ttl),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	filePath := filepath.Join(c.dir, name+".json")
	return os.WriteFile(filePath, data, 0644)
}

func (c *Cache) CleanExpired() error {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	now := time.Now().Unix()

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(c.dir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var cacheEntry cacheEntry
		if err := json.Unmarshal(data, &cacheEntry); err != nil {
			continue
		}

		if now > cacheEntry.ExpiresAt {
			os.Remove(filePath)
		}
	}

	return nil
}

// GenerateCacheKey generates a cache key from cwd and action name
// Format: {projectName}_{parentHashFirst4Chars}_{actionName}
func (c *Cache) GenerateCacheKey(cwd string, actionName string) string {
	projectName := filepath.Base(cwd)
	parentPath := filepath.Dir(cwd)
	
	// Generate hash of parent path
	hash := sha256.Sum256([]byte(parentPath))
	hashStr := fmt.Sprintf("%x", hash[:2]) // First 2 bytes = 4 hex chars
	
	return fmt.Sprintf("%s_%s_%s", projectName, hashStr, actionName)
}

// GetWithCwd retrieves a cached value using cwd and action name
func (c *Cache) GetWithCwd(cwd string, actionName string) (string, bool) {
	cacheKey := c.GenerateCacheKey(cwd, actionName)
	return c.Get(cacheKey)
}

// SetWithCwd stores a value in cache using cwd and action name
func (c *Cache) SetWithCwd(cwd string, actionName string, result string, ttl int) error {
	cacheKey := c.GenerateCacheKey(cwd, actionName)
	return c.Set(cacheKey, result, ttl)
}

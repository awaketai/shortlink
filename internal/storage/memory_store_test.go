package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMemoryStore_FindByShortCode(t *testing.T) {
	// 使用真实的 MemoryStore 实例（遵循宪法 2.3：拒绝 Mocks）
	seed := func(t *testing.T) *MemoryStore {
		t.Helper()
		store := NewMemoryStore()
		// 预置测试数据
		links := []Link{
			{ShortCode: "abc123", LongURL: "https://example.com"},
			{ShortCode: "xyz789", LongURL: "https://golang.org"},
		}
		for _, l := range links {
			if err := store.Save(context.Background(), l); err != nil {
				t.Fatalf("seed data failed: %v", err)
			}
		}
		return store
	}

	// 表格驱动测试（遵循宪法 2.2）
	tests := []struct {
		name      string
		shortCode string
		wantURL   string
		wantErr   error
	}{
		{
			name:      "existing short code returns correct link",
			shortCode: "abc123",
			wantURL:   "https://example.com",
			wantErr:   nil,
		},
		{
			name:      "another existing short code returns correct link",
			shortCode: "xyz789",
			wantURL:   "https://golang.org",
			wantErr:   nil,
		},
		{
			name:      "non-existent short code returns ErrNotFound",
			shortCode: "nonexistent",
			wantURL:   "",
			wantErr:   ErrNotFound,
		},
		{
			name:      "empty short code returns ErrNotFound",
			shortCode: "",
			wantURL:   "",
			wantErr:   ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := seed(t)

			got, err := store.FindByShortCode(context.Background(), tt.shortCode)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindByShortCode(%q) error = %v, wantErr = %v", tt.shortCode, err, tt.wantErr)
				return
			}

			if tt.wantErr != nil {
				if got != nil {
					t.Errorf("FindByShortCode(%q) returned non-nil link on error", tt.shortCode)
				}
				return
			}

			if got == nil {
				t.Fatalf("FindByShortCode(%q) returned nil link, want non-nil", tt.shortCode)
			}
			if got.LongURL != tt.wantURL {
				t.Errorf("FindByShortCode(%q).LongURL = %q, want %q", tt.shortCode, got.LongURL, tt.wantURL)
			}
			if got.ShortCode != tt.shortCode {
				t.Errorf("FindByShortCode(%q).ShortCode = %q, want %q", tt.shortCode, got.ShortCode, tt.shortCode)
			}
		})
	}
}

func TestMemoryStore_Save(t *testing.T) {
	// 表格驱动测试（遵循宪法 2.2）
	tests := []struct {
		name     string
		link     Link
		preSeed  []Link // 预先植入的数据
		wantErr  error
		setupErr error // 预期设置阶段的错误
	}{
		{
			name:    "successful save of new link",
			link:    Link{ShortCode: "abc123", LongURL: "https://example.com"},
			wantErr: nil,
		},
		{
			name:    "successful save of another new link",
			link:    Link{ShortCode: "xyz789", LongURL: "https://golang.org"},
			wantErr: nil,
		},
		{
			name:    "save with empty short code should succeed",
			link:    Link{ShortCode: "", LongURL: "https://empty.com"},
			wantErr: nil,
		},
		{
			name:    "save with empty long URL should succeed",
			link:    Link{ShortCode: "emptyurl", LongURL: ""},
			wantErr: nil,
		},
		{
			name:    "save with special characters in short code should succeed",
			link:    Link{ShortCode: "abc-123_456", LongURL: "https://special-chars.com"},
			wantErr: nil,
		},
		{
			name:    "save with very long URL should succeed",
			link:    Link{ShortCode: "long", LongURL: "https://example.com/very/long/path/with/many/segments/and/parameters?param1=value1&param2=value2&param3=value3"},
			wantErr: nil,
		},
		{
			name: "save with existing short code should return ErrShortCodeExists",
			link: Link{ShortCode: "existing", LongURL: "https://new.com"},
			preSeed: []Link{
				{ShortCode: "existing", LongURL: "https://existing.com"},
			},
			wantErr: ErrShortCodeExists,
		},
		{
			name: "save with same short code but different URL should return ErrShortCodeExists",
			link: Link{ShortCode: "duplicate", LongURL: "https://different.com"},
			preSeed: []Link{
				{ShortCode: "duplicate", LongURL: "https://original.com"},
			},
			wantErr: ErrShortCodeExists,
		},
		{
			name: "save with same URL but different short code should succeed",
			link: Link{ShortCode: "newcode", LongURL: "https://same-url.com"},
			preSeed: []Link{
				{ShortCode: "oldcode", LongURL: "https://same-url.com"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用真实的 MemoryStore 实例（遵循宪法 2.3：拒绝 Mocks）
			store := NewMemoryStore()

			// 设置预置数据
			for _, link := range tt.preSeed {
				if err := store.Save(context.Background(), link); err != nil {
					t.Fatalf("seed data failed: %v", err)
				}
			}

			// 记录保存前的时间戳
			beforeSave := time.Now()

			// 执行被测试的操作
			err := store.Save(context.Background(), tt.link)

			// 验证错误
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Save() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			// 如果预期有错误，验证原有数据未被覆盖
			if tt.wantErr != nil {
				// 如果有预置数据，验证原有数据未被覆盖
				if len(tt.preSeed) > 0 {
					found, err := store.FindByShortCode(context.Background(), tt.link.ShortCode)
					if err != nil {
						t.Errorf("Save() error case: failed to find existing link: %v", err)
						return
					}
					if found.LongURL != tt.preSeed[0].LongURL {
						t.Errorf("Save() should not have overwritten existing link, got URL %q, want %q",
							found.LongURL, tt.preSeed[0].LongURL)
					}
				}
				return
			}

			// 验证保存成功的情况
			if err != nil {
				t.Errorf("Save() unexpected error = %v", err)
				return
			}

			// 验证数据确实被保存
			saved, err := store.FindByShortCode(context.Background(), tt.link.ShortCode)
			if err != nil {
				t.Errorf("Save() failed to find saved link: %v", err)
				return
			}

			if saved == nil {
				t.Fatalf("Save() returned nil link after successful save")
			}

			// 验证保存的数据正确性
			if saved.ShortCode != tt.link.ShortCode {
				t.Errorf("Save().ShortCode = %q, want %q", saved.ShortCode, tt.link.ShortCode)
			}
			if saved.LongURL != tt.link.LongURL {
				t.Errorf("Save().LongURL = %q, want %q", saved.LongURL, tt.link.LongURL)
			}

			// 验证 CreatedAt 字段被正确设置
			if saved.CreatedAt.Before(beforeSave) {
				t.Errorf("Save().CreatedAt = %v, should be after %v", saved.CreatedAt, beforeSave)
			}
			if saved.CreatedAt.After(time.Now().Add(time.Second)) {
				t.Errorf("Save().CreatedAt = %v, should not be too far in the future", saved.CreatedAt)
			}
		})
	}
}

func TestMemoryStore_Save_Concurrent(t *testing.T) {
	// 测试并发安全性
	store := NewMemoryStore()

	const numGoroutines = 100
	const numSaves = 10

	var wg sync.WaitGroup
	errChan := make(chan error, numGoroutines*numSaves)

	// 启动多个 goroutine 并发保存不同的链接
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numSaves; j++ {
				link := Link{
					ShortCode: fmt.Sprintf("link_%d_%d", id, j),
					LongURL:   fmt.Sprintf("https://example.com/%d/%d", id, j),
				}
				if err := store.Save(context.Background(), link); err != nil {
					errChan <- fmt.Errorf("goroutine %d, save %d: %w", id, j, err)
				}
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		t.Error(err)
	}

	// 验证所有链接都被正确保存
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numSaves; j++ {
			shortCode := fmt.Sprintf("link_%d_%d", i, j)
			link, err := store.FindByShortCode(context.Background(), shortCode)
			if err != nil {
				t.Errorf("Concurrent save failed to find link %s: %v", shortCode, err)
			}
			if link == nil {
				t.Errorf("Concurrent save: link %s is nil", shortCode)
			}
			expectedURL := fmt.Sprintf("https://example.com/%d/%d", i, j)
			if link != nil && link.LongURL != expectedURL {
				t.Errorf("Concurrent save: link %s has URL %q, want %q", shortCode, link.LongURL, expectedURL)
			}
		}
	}
}

func TestMemoryStore_Save_ConcurrentSameShortCode(t *testing.T) {
	// 测试并发保存相同 short code 的行为
	store := NewMemoryStore()

	const numGoroutines = 50
	shortCode := "conflict"
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// 启动多个 goroutine 尝试保存相同的 short code
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			link := Link{
				ShortCode: shortCode,
				LongURL:   fmt.Sprintf("https://example.com/%d", id),
			}
			err := store.Save(context.Background(), link)
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// 应该只有一个保存成功
	if successCount != 1 {
		t.Errorf("Concurrent save with same short code: expected 1 success, got %d", successCount)
	}

	// 验证确实只有一个链接被保存
	saved, err := store.FindByShortCode(context.Background(), shortCode)
	if err != nil {
		t.Errorf("Concurrent save: failed to find the saved link: %v", err)
	}
	if saved == nil {
		t.Error("Concurrent save: no link was saved")
	}
}

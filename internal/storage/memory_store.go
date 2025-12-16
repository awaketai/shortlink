package storage

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Link 现在属于storage包，如果shortener.service要操作link
// 就需要引入storage包并使用storage.Link，这会产生一个从业务核心到具体存储实现的依赖
// 不太理想，这是后续接口抽象要解决的问题
type Link struct {
	ShortCode  string
	LongURL    string
	VisitCount int64
	CreatedAt  time.Time
}

var ErrNotFound = errors.New("storage: link not found")
var ErrShortCodeExists = errors.New("storage: short code already exists")

type Store struct {
	mu    sync.RWMutex
	links map[string]Link
}

func NewStore() *Store {
	return &Store{
		links: make(map[string]Link),
	}
}

func (s *Store) Save(ctx context.Context, link Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.links[link.ShortCode]; ok {
		return ErrShortCodeExists
	}

	link.CreatedAt = time.Now()
	s.links[link.ShortCode] = link

	return nil
}

func (s *Store) FindByShortCode(ctx context.Context, shortCode string) (Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if link, ok := s.links[shortCode]; ok {
		return link, nil
	}

	return Link{}, ErrNotFound
}

func (s *Store) IncrementVisitCount(ctx context.Context, shortCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if link, ok := s.links[shortCode]; ok {
		link.VisitCount++
		s.links[shortCode] = link
		return nil
	}

	return ErrNotFound
}

func (log *Store) Close() error {
	return nil
}

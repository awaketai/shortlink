package storage

import (
	"context"
	"sync"
	"time"
)

type MemoryStore struct {
	mu    sync.RWMutex
	links map[string]*Link
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		links: make(map[string]*Link),
	}
}

func (s *MemoryStore) Save(ctx context.Context, link Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.links[link.ShortCode]; ok {
		return ErrShortCodeExists
	}

	link.CreatedAt = time.Now()
	s.links[link.ShortCode] = &link

	return nil
}

func (s *MemoryStore) FindByShortCode(ctx context.Context, shortCode string) (*Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if link, ok := s.links[shortCode]; ok {
		return link, nil
	}

	return nil, ErrNotFound
}

func (s *MemoryStore) IncrementVisitCount(ctx context.Context, shortCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if link, ok := s.links[shortCode]; ok {
		link.VisitCount++
		s.links[shortCode] = link
		return nil
	}

	return ErrNotFound
}

func (log *MemoryStore) Close() error {
	return nil
}

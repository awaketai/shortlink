package storage

import (
	"context"
	"sync"
	"time"
)

type Store struct {
	mu    sync.RWMutex
	links map[string]*Link
}

func NewStore() *Store {
	return &Store{
		links: make(map[string]*Link),
	}
}

func (s *Store) Save(ctx context.Context, link Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.links[link.ShortCode]; ok {
		return ErrShortCodeExists
	}

	link.CreatedAt = time.Now()
	s.links[link.ShortCode] = &link

	return nil
}

func (s *Store) FindByShortCode(ctx context.Context, shortCode string) (*Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if link, ok := s.links[shortCode]; ok {
		return link, nil
	}

	return nil, ErrNotFound
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

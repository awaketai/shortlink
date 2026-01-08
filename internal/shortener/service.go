package shortener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"shortlink/internal/idgen"
	"shortlink/internal/storage"
	"strings"
	"time"
)

// 核心业务逻辑
// 职责：实现创建短链接，获取场链接等核心业务流程，编排对idgen.Generator和storage.Store的调用
// 内聚性：聚焦核心短链接业务逻辑
//
// SOLID原则
// SRP(单一职责)：shortener包专注于业务编排，将ID生成和存储的具体实现委托给其他包
// DIP(依赖倒置原则)：

var (
	ErrInvalidLongURL            = errors.New("shortener: long URL is invalid or empty.")
	ErrShortCodeTooShort         = errors.New("shortener: short code is too short or invalid.")
	ErrShortCodeGenerationFailed = errors.New("shortener: failed to generate short code.")
	ErrLinkNotFound              = errors.New("shortener: link not found.")
	ErrConflict                  = errors.New("shortener: conflict,possibly short code exists or generation failed after retries.")
)

type Config struct {
	Store           storage.Storer
	Generator       idgen.Generator
	Logger          *log.Logger
	MaxGenAttemps   int
	MinShortCodeLen int
}

type Service struct {
	store           storage.Storer
	generator       idgen.Generator
	logger          *log.Logger
	maxGenAttempts  int
	minShortCodeLen int
}

func NewService(cfg Config) *Service {
	if cfg.Store == nil || cfg.Generator == nil {
		log.Fatal("storage and generator are required")
	}
	if cfg.MaxGenAttemps <= 0 {
		cfg.MaxGenAttemps = 3
	}
	if cfg.MinShortCodeLen <= 0 {
		cfg.MinShortCodeLen = 5
	}

	return &Service{
		store:           cfg.Store,
		generator:       cfg.Generator,
		logger:          cfg.Logger,
		maxGenAttempts:  cfg.MaxGenAttemps,
		minShortCodeLen: cfg.MinShortCodeLen,
	}
}

func (s *Service) CreateShortLink(ctx context.Context, longURL string) (string, error) {
	if strings.TrimSpace(longURL) == "" {
		return "", ErrInvalidLongURL
	}
	var shortCode string
	for i := range s.maxGenAttempts {
		log.Printf("DEBUG: Attempting to generate short code,attempt %d,longURL: %s \n", i+1, longURL)
		code, genErr := s.generator.GenerateShortCode(ctx, longURL)
		if genErr != nil {
			return "", fmt.Errorf("attempt %d:failed to generate short code:%w", i+1, genErr)
		}
		shortCode = code
		if len(shortCode) < s.minShortCodeLen {
			s.logger.Printf("WARN: Generated short code too short, retrying. Code: %s, Attempt: %d\n", shortCode, i+1)
			if i < s.maxGenAttempts {
				continue
			} else {
				break
			}
		}

		linkToSave := storage.Link{
			ShortCode:  shortCode,
			LongURL:    longURL,
			VisitCount: 0,
			CreatedAt:  time.Now().UTC(),
		}
		saveErr := s.store.Save(ctx, linkToSave)
		if saveErr != nil {
			if errors.Is(saveErr, storage.ErrShortCodeExists) && i < s.maxGenAttempts-1 {
				log.Printf("WARN: Short code collision,retrying,Attempt: %d Code:%s\n", i+1, shortCode)
				continue
			}
			return "", fmt.Errorf("attempt %d:failed to save short link:%w", i+1, saveErr)
		}

		return shortCode, nil
	}

	return "", fmt.Errorf("failed to generate short code after %d attempts", s.maxGenAttempts)
}

func (s *Service) GetAndTrackLongURL(ctx context.Context, shortCode string) (string, error) {
	if len(shortCode) < s.minShortCodeLen {
		return "", ErrShortCodeTooShort
	}
	link, err := s.store.FindByShortCode(ctx, shortCode)
	if err != nil {
		s.logger.Printf("INFO: Short code not found in store. ShortCode: %s\n", shortCode)
		if errors.Is(err, storage.ErrNotFound) {
			return "", fmt.Errorf("for code '%s': %w", shortCode, ErrLinkNotFound)
		}
	}

	go func(sc string, currentCount int64) {
		bgCtx := context.Background()
		if err := s.store.IncrementVisitCount(bgCtx, sc); err != nil {
			log.Printf("ERROR: Failed to increment visit count (async).ShortCode: %s,Error:%v", sc, err)
		}
		log.Printf("INFO: Visit count incremented successfully.ShortCode: %s,CurrentCount:%d", sc, currentCount+1)
	}(shortCode, link.VisitCount)

	return link.LongURL, nil
}

func preview(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}

	return s
}

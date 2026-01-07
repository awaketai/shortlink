package storage

import (
	"context"
	"errors"
	"time"
)

// Link 现在属于storage包，如果shortener.service要操作link
// 就需要引入storage包并使用storage.Link，这会产生一个从业务核心到具体存储实现的依赖
type Link struct {
	ShortCode string
	LongURL   string
	// 访问次数，处理层需要处理并发更新问题
	VisitCount int64
	CreatedAt  time.Time
}

// Storer 数据存储层需要提供的核心能力
// 所有实现都应该是并发安全的
type Storer interface {
	// Save 保存一个新的短链接映射，如果 shortCode已经存在，则返回 ErrShortCodeExists
	// 如果 longURL 无效，实现可以返回特定错误或依赖上层校验
	Save(ctx context.Context, link Link) error
	// FindByShortCode 根据短码查找对应的 Link 信息，如果未找到，返回 ErrNotFound
	FindByShortCode(ctx context.Context, shortCode string) (*Link, error)
	// IncrementVisitCount 原子增加短码的访问次数，如果shortCode不存在，可以返回 ErrNotFound，或者静默失败，取决于具体业务需求
	// 此方法必须是并发安全的
	IncrementVisitCount(ctx context.Context, shortCode string) error
	// Close 关闭并释放存储层占用的资源(如果数据库连接池)，应确保幂等性，多次调用 Close 不会产生副作用
	Close() error
}

var (
	ErrNotFound        = errors.New("storage: link not found")
	ErrShortCodeExists = errors.New("storage: short code already exists")
)

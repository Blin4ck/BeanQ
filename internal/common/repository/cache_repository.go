package repository

import (
	"context"
	"time"
)

// CacheRepository определяет методы для работы с кэшем.
type CacheRepository interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

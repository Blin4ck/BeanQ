package repository

import (
	"context"
	"time"
)

// TokenRepository определяет методы для хранения и получения токенов пользователя.
type TokenRepository interface {
	SetToken(ctx context.Context, userID string, token string, ttl time.Duration) error
	GetToken(ctx context.Context, userID string) (string, error)
	DeleteToken(ctx context.Context, userID string) error
}

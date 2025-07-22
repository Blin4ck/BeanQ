package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) *TokenRepository {
	return &TokenRepository{client: client}
}

func (r *TokenRepository) SetToken(ctx context.Context, userID string, token string, ttl time.Duration) error {
	return r.client.Set(ctx, r.key(userID), token, ttl).Err()
}

func (r *TokenRepository) key(userID string) string {
	return "refresh_token:" + userID
}

func (r *TokenRepository) GetToken(ctx context.Context, userID string) (string, error) {
	return r.client.Get(ctx, r.key(userID)).Result()
}

func (r *TokenRepository) DeleteToken(ctx context.Context, userID string) error {
	return r.client.Del(ctx, r.key(userID)).Err()
}

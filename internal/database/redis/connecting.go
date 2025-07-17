package redis

import (
	"github.com/redis/go-redis/v9"
)

// Config содержит параметры подключения к Redis.
type Config struct {
	Addr     string
	Password string
	DB       int
}

// NewRedisClient создает и возвращает новый клиент Redis.
func NewRedisClient(cfg Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

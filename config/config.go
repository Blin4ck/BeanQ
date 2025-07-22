package config

import (
	"os"
	"strconv"
)

// Config содержит параметры подключения к базе данных.
type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBMaxIdleConns    int
	DBMaxOpenConns    int
	DBConnMaxLifetime int

	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

// New создает новый экземпляр Config, заполняя его из переменных окружения.
func New() *Config {
	return &Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "password"),
		DBName:            getEnv("DB_NAME", "coffee_shop"),
		DBMaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
		DBMaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
		DBConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 3600),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt возвращает значение переменной окружения как int или значение по умолчанию.
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

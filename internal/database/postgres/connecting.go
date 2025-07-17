package postgres

import (
	"coffe/config"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect устанавливает соединение с базой данных PostgreSQL.
func Connect(config *config.Config) (*gorm.DB, error) {

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	gormLoger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	})

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: gormLoger,
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка подключения: %w", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("ошибка получения sqlDB: %w", err)
	}

	sqlDB.SetMaxIdleConns(config.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.DBConnMaxLifetime) * time.Second)

	return db, nil
}

// Close закрывает соединение с базой данных.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()

	if err != nil {
		return fmt.Errorf("ошибка получения sqlDB: %w", err)
	}

	return sqlDB.Close()
}

// AutoMigrate выполняет автоматическую миграцию моделей в базе данных.
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}

package repository

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func RetryConectDB(dsn string, maxRetries int) (*gorm.DB, error) {

	var db *gorm.DB
	var err error

	delay := time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}

		time.Sleep(delay)
		delay *= 2
	}

	return nil, err

}

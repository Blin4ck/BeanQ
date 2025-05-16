package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func ConectDB(env string) (*gorm.DB, error) {

	fmt.Print(env)
	db, err := gorm.Open(postgres.Open(env), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return db, nil
}

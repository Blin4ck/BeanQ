package repository

import (
	"test/models"
)

func (r *Repository) Migrate() error {
	err := r.DB.AutoMigrate(
		&models.Client{},
		&models.Position{},
		&models.Menu{},
		&models.Order{},
	)
	if err != nil {
		return err
	}
	return nil
}

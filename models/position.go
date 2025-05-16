package models

import "gorm.io/gorm"

type Position struct {
	gorm.Model
	ID          uint        `gorm:"primaryKey"`
	Name        string      `json:"name"`
	Price       int         `json:"price"`
	Description string      `json:"description"`
	Category    string      `json:"category"`
	OrderItems  []OrderItem `gorm:"foreignKey:PositionID"`
}

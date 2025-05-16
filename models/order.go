package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID         uint        `gorm:"primaryKey"`
	Client     Client      `gorm:"foreignKey:ClientID" json:"-"`
	CreatedAt  time.Time   `json:"created_at"`
	Price      int         `json:"price"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
}

type OrderItem struct {
	gorm.Model
	ID         uint     `gorm:"primaryKey"`
	OrderID    uint     `json:"order_id"`
	PositionID uint     `json:"position_id"`
	Position   Position `gorm:"foreignKey:PositionID" json:"position"`
	Quantity   int      `json:"quantity"`
	Price      int      `json:"price"`
}

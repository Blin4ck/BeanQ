package models

import (
	"time"
)

type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	ClientID   uint        `json:"client_id"`
	Client     Client      `gorm:"foreignKey:ClientID" json:"-"`
	CreatedAt  time.Time   `json:"created_at"`
	Price      int         `json:"price"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
}

type OrderItem struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	OrderID    uint     `json:"order_id"`
	PositionID uint     `json:"position_id"`
	Position   Position `gorm:"foreignKey:PositionID" json:"position"`
	Quantity   int      `json:"quantity"`
	Price      int      `json:"price"`
}

package dto

import (
	"time"

	"github.com/google/uuid"
)

type MenuSearchDTO struct {
	Query       string     `json:"query"`
	MenuID      uuid.UUID  `json:"menu_id"`
	CategoryID  uuid.UUID  `json:"category_id"`
	ProductID   uuid.UUID  `json:"product_id"`
	IsActive    *bool      `json:"is_active"`
	PriceRange  [2]float64 `json:"price_range"`
	ValidAfter  *time.Time `json:"valid_after"`
	ValidBefore *time.Time `json:"valid_before"`
	Pagination  Pagination `json:"pagination"`
	Sorting     Sorting    `json:"sorting"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type Sorting struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

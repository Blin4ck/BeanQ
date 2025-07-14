package entity

import (
	"time"

	"github.com/google/uuid"
)

// MenuCategory представляет категорию в меню.
type MenuCategory struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`               // "Кофе", "Чай", "Десерты"
	Description string    `json:"description" db:"description"` // описание категории
	SortOrder   int       `json:"sort_order" db:"sort_order"`   // порядок отображения
	IsActive    bool      `json:"is_active" db:"is_active"`     // активна ли категория
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Menu представляет меню кофейни.
type Menu struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`               // "Основное меню", "Сезонное меню"
	Description string         `json:"description" db:"description"` // описание меню
	Categories  []MenuCategory `json:"categories" db:"categories"`   // категории в меню
	IsActive    bool           `json:"is_active" db:"is_active"`     // активно ли меню
	ValidFrom   time.Time      `json:"valid_from" db:"valid_from"`   // с какого времени действует
	ValidTo     *time.Time     `json:"valid_to" db:"valid_to"`       // до какого времени действует (nil = бессрочно)
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// MenuItem представляет позицию в меню.
type MenuItem struct {
	ID         uuid.UUID     `json:"id" db:"id"`
	MenuID     uuid.UUID     `json:"menu_id" db:"menu_id"`
	ProductID  uuid.UUID     `json:"product_id" db:"product_id"`
	Product    *Product      `json:"product,omitempty" db:"product"`
	CategoryID uuid.UUID     `json:"category_id" db:"category_id"`
	Category   *MenuCategory `json:"category,omitempty" db:"category"`
	SortOrder  int           `json:"sort_order" db:"sort_order"` // порядок в категории
	IsActive   bool          `json:"is_active" db:"is_active"`   // активна ли позиция
	CreatedAt  time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" db:"updated_at"`
}

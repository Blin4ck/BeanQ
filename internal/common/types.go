package common

import (
	"github.com/google/uuid"
)

// User представляет пользователя системы.
type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Surname  string    `json:"surname" db:"surname"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"-" db:"password"`
	Role     *Role     `json:"role,omitempty" db:"role"`
}

// Role представляет роль пользователя.
type Role struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}

// Product представляет продукт меню.
type Product struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Category    string    `json:"category" db:"category"`
	Price       float64   `json:"price" db:"price"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	ImageURL    string    `json:"image_url" db:"image_url"`
}

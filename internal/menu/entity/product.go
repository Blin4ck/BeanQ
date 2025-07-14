package entity

import (
	"time"

	"github.com/google/uuid"
)

// Product представляет продукт меню.
type Product struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	Name        string        `json:"name" db:"name"`
	Category    string        `json:"category" db:"category"` // "coffee", "dessert"
	Price       float64       `json:"price" db:"price"`
	Description string        `json:"description" db:"description"`
	IsActive    bool          `json:"is_active" db:"is_active"`
	ImageURL    string        `json:"image_url" db:"image_url"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
	Ingredients []*Ingredient `json:"ingredients,omitempty" db:"ingredients"` // ингредиенты продукта
}

// Ingredient представляет ингредиент продукта.
type Ingredient struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Quantity float64   `json:"quantity" db:"quantity"`
	Unit     string    `json:"unit" db:"unit"` // "ml", "g", "шт"
}

// ProductIngredient представляет связь между продуктом и ингредиентом.
type ProductIngredient struct {
	ProductID    uuid.UUID   `json:"product_id" db:"product_id"`
	IngredientID uuid.UUID   `json:"ingredient_id" db:"ingredient_id"`
	Quantity     float64     `json:"quantity" db:"quantity"` // количество ингредиента в продукте
	Unit         string      `json:"unit" db:"unit"`         // единица измерения для этого продукта
	Product      *Product    `json:"product,omitempty" db:"product"`
	Ingredient   *Ingredient `json:"ingredient,omitempty" db:"ingredient"`
}

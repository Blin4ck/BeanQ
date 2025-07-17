package repository

import (
	"coffe/internal/menu/entity"
	"context"

	"github.com/google/uuid"
)

// IngredientRepository определяет методы для работы с ингредиентами.
type IngredientRepository interface {
	Create(ctx context.Context, ingredient *entity.Ingredient) error
	GetByName(ctx context.Context, name string) (*entity.Ingredient, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Ingredient, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context) ([]*entity.Ingredient, error)
	Update(ctx context.Context, ingredient *entity.Ingredient) error
	Search(ctx context.Context, query string) ([]*entity.Ingredient, error)
	GetByUnit(ctx context.Context, unit string) ([]*entity.Ingredient, error)
	GetByQuantityRange(ctx context.Context, min, max float64) ([]*entity.Ingredient, error)
	GetIngredientsByProduct(ctx context.Context, productID uuid.UUID) ([]*entity.Ingredient, error)
	AddIngredientToProduct(ctx context.Context, productID, ingredientID uuid.UUID, quantity float64) error
	RemoveIngredientFromProduct(ctx context.Context, productID, ingredientID uuid.UUID) error
	GetProductsByIngredient(ctx context.Context, ingredientID uuid.UUID) ([]*entity.Product, error)
	ValidateIngredient(ctx context.Context, ingredient *entity.Ingredient) error                   // валидация ингредиента
	CheckAvailability(ctx context.Context, ingredientID uuid.UUID, quantity float64) (bool, error) // проверка доступности
	Count(ctx context.Context) (int64, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

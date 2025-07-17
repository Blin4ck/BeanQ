package repository

import (
	"coffe/internal/menu/entity"
	"context"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error                                  // создание товара
	CreateIngredient(ctx context.Context, ingredient *entity.Ingredient) error                  // создание ингредиента
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)                         // поиск по id
	Update(ctx context.Context, product *entity.Product) error                                  // обновление товара
	Delete(ctx context.Context, id uuid.UUID) error                                             // удаление товара
	GetByCategory(ctx context.Context, category string) ([]*entity.Product, error)              // товары по категории
	GetActive(ctx context.Context) ([]*entity.Product, error)                                   // только активные товары
	Count(ctx context.Context) (int64, error)                                                   // количество товаров
	GetNotActive(ctx context.Context) ([]*entity.Product, error)                                // неактивные товары
	GetAll(ctx context.Context) ([]*entity.Product, error)                                      // все товары
	Search(ctx context.Context, query string) ([]*entity.Product, error)                        // поиск по названию/описанию
	GetByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*entity.Product, error) // товары по цене
	Activate(ctx context.Context, id uuid.UUID) error                                           // активировать товар
	Deactivate(ctx context.Context, id uuid.UUID) error                                         // деактивировать товар
	Exists(ctx context.Context, id uuid.UUID) (bool, error)                                     // проверить существование
	// Методы для работы с ингредиентами
	GetIngredientsByProduct(ctx context.Context, productID uuid.UUID) ([]*entity.Ingredient, error)                      // получить ингредиенты продукта
	AddIngredientToProduct(ctx context.Context, productID, ingredientID uuid.UUID, quantity float64, unit string) error  // добавить ингредиент к продукту
	RemoveIngredientFromProduct(ctx context.Context, productID, ingredientID uuid.UUID) error                            // удалить ингредиент из продукта
	UpdateProductIngredient(ctx context.Context, productID, ingredientID uuid.UUID, quantity float64, unit string) error // обновить количество ингредиента
}

package usecase

import (
	"coffe/internal/menu/entity"
	"coffe/internal/menu/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

// ProductUsecase реализует бизнес-логику для работы с продуктами.
type ProductUsecase struct {
	productRepo repository.ProductRepository
}

// NewProductUsecase создает новый экземпляр ProductUsecase.
func NewProductUsecase(productRepo repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{productRepo: productRepo}
}

// Create добавляет новый продукт.
func (u *ProductUsecase) Create(ctx context.Context, product *entity.Product) error {
	if product.Name == "" {
		return errors.New("имя не может быть пустым")
	}
	if product.Category == "" {
		return errors.New("вы не задали категорию")
	}
	if product.Price <= 0 {
		return errors.New("цена не может быть отрицательной")
	}
	if product.Description == "" {
		return errors.New("вы не задали описание")
	}

	return u.productRepo.Create(ctx, product)
}

// CreateIngredient добавляет новый ингредиент.
func (u *ProductUsecase) CreateIngredient(ctx context.Context, ingredient *entity.Ingredient) error {

	if ingredient.Name == "" {
		return errors.New("имя не может быть пустым")
	}
	if ingredient.Quantity <= 0 {
		return errors.New("количество не может быть отрицательным")
	}
	return u.productRepo.CreateIngredient(ctx, ingredient)
}

// GetByID возвращает продукт по идентификатору.
func (u *ProductUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	if id == uuid.Nil {
		return nil, errors.New("id не может быть пустым")
	}
	return u.productRepo.GetByID(ctx, id)
}

// Update обновляет продукт.
func (u *ProductUsecase) Update(ctx context.Context, product *entity.Product) error {
	if product.ID == uuid.Nil {
		return errors.New("id не может быть пустым")
	}
	return u.productRepo.Update(ctx, product)
}

// Delete удаляет продукт по идентификатору.
func (u *ProductUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id не может быть пустым")
	}
	return u.productRepo.Delete(ctx, id)
}

// GetByCategory возвращает продукты по категории.
func (u *ProductUsecase) GetByCategory(ctx context.Context, category string) ([]*entity.Product, error) {
	if category == "" {
		return nil, errors.New("вы не задали категорию")
	}
	return u.productRepo.GetByCategory(ctx, category)
}

// GetActive возвращает активные продукты.
func (u *ProductUsecase) GetActive(ctx context.Context) ([]*entity.Product, error) {
	return u.productRepo.GetActive(ctx)
}

// Count возвращает количество продуктов.
func (u *ProductUsecase) Count(ctx context.Context) (int64, error) {
	return u.productRepo.Count(ctx)
}

// GetNotActive возвращает неактивные продукты.
func (u *ProductUsecase) GetNotActive(ctx context.Context) ([]*entity.Product, error) {
	return u.productRepo.GetNotActive(ctx)
}

// GetAll возвращает все продукты.
func (u *ProductUsecase) GetAll(ctx context.Context) ([]*entity.Product, error) {
	products, err := u.productRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.New("ошибка при получении продуктов")
	}
	return products, nil
}

// Search ищет продукты по запросу.
func (u *ProductUsecase) Search(ctx context.Context, query string) ([]*entity.Product, error) {
	if query == "" {
		return nil, errors.New("поисковый запрос не может быть пустым")
	}

	products, err := u.productRepo.Search(ctx, query)
	if err != nil {
		return nil, errors.New("ошибка при поиске продуктов")
	}
	return products, nil
}

// GetByPriceRange возвращает продукты в заданном диапазоне цен.
func (u *ProductUsecase) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*entity.Product, error) {
	if minPrice < 0 || maxPrice < 0 {
		return nil, errors.New("цена не может быть отрицательной")
	}
	if minPrice > maxPrice {
		return nil, errors.New("минимальная цена не может быть больше максимальной")
	}

	products, err := u.productRepo.GetByPriceRange(ctx, minPrice, maxPrice)
	if err != nil {
		return nil, errors.New("ошибка при получении продуктов по диапазону цен")
	}
	return products, nil
}

func (u *ProductUsecase) Activate(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}

	// Проверяем, существует ли продукт
	_, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("продукт не найден")
	}

	return u.productRepo.Activate(ctx, id)
}

func (u *ProductUsecase) Deactivate(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}

	// Проверяем, существует ли продукт
	_, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("продукт не найден")
	}

	return u.productRepo.Deactivate(ctx, id)
}

func (u *ProductUsecase) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, errors.New("ID не может быть пустым")
	}

	exists, err := u.productRepo.Exists(ctx, id)
	if err != nil {
		return false, errors.New("ошибка при проверке существования продукта")
	}

	return exists, nil
}

// validateProduct валидирует продукт
func (u *ProductUsecase) validateProduct(product *entity.Product) error {
	if product == nil {
		return errors.New("продукт не может быть пустым")
	}
	if product.Name == "" {
		return errors.New("название продукта не может быть пустым")
	}
	if product.Category == "" {
		return errors.New("категория продукта не может быть пустой")
	}
	if product.Price <= 0 {
		return errors.New("цена продукта должна быть больше нуля")
	}
	if product.Description == "" {
		return errors.New("описание продукта не может быть пустым")
	}

	return nil
}

// GetProductIngredients получает ингредиенты продукта
func (u *ProductUsecase) GetProductIngredients(ctx context.Context, productID uuid.UUID) ([]*entity.Ingredient, error) {
	if productID == uuid.Nil {
		return nil, errors.New("ID продукта не может быть пустым")
	}

	// Проверяем, существует ли продукт
	_, err := u.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("продукт не найден")
	}

	ingredients, err := u.productRepo.GetIngredientsByProduct(ctx, productID)
	if err != nil {
		return nil, errors.New("ошибка при получении ингредиентов продукта")
	}

	return ingredients, nil
}

// AddIngredientToProduct добавляет ингредиент к продукту
func (u *ProductUsecase) AddIngredientToProduct(ctx context.Context, productID, ingredientID uuid.UUID, quantity float64, unit string) error {
	if productID == uuid.Nil {
		return errors.New("ID продукта не может быть пустым")
	}
	if ingredientID == uuid.Nil {
		return errors.New("ID ингредиента не может быть пустым")
	}
	if quantity <= 0 {
		return errors.New("количество должно быть больше нуля")
	}
	if unit == "" {
		return errors.New("единица измерения не может быть пустой")
	}

	return u.productRepo.AddIngredientToProduct(ctx, productID, ingredientID, quantity, unit)
}

// RemoveIngredientFromProduct удаляет ингредиент из продукта
func (u *ProductUsecase) RemoveIngredientFromProduct(ctx context.Context, productID, ingredientID uuid.UUID) error {
	if productID == uuid.Nil {
		return errors.New("ID продукта не может быть пустым")
	}
	if ingredientID == uuid.Nil {
		return errors.New("ID ингредиента не может быть пустым")
	}

	return u.productRepo.RemoveIngredientFromProduct(ctx, productID, ingredientID)
}

// UpdateProductIngredient обновляет количество ингредиента в продукте
func (u *ProductUsecase) UpdateProductIngredient(ctx context.Context, productID, ingredientID uuid.UUID, quantity float64, unit string) error {
	if productID == uuid.Nil {
		return errors.New("ID продукта не может быть пустым")
	}
	if ingredientID == uuid.Nil {
		return errors.New("ID ингредиента не может быть пустым")
	}
	if quantity <= 0 {
		return errors.New("количество должно быть больше нуля")
	}
	if unit == "" {
		return errors.New("единица измерения не может быть пустой")
	}

	return u.productRepo.UpdateProductIngredient(ctx, productID, ingredientID, quantity, unit)
}

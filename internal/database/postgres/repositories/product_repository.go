package repositories

import (
	"coffe/internal/menu/entity"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create создает новый продукт
func (r *ProductRepository) Create(ctx context.Context, product *entity.Product) error {
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

	return r.db.WithContext(ctx).Create(product).Error
}

// CreateIngredient создает ингредиент (должно быть в IngredientRepository)
func (r *ProductRepository) CreateIngredient(ctx context.Context, ingredient *entity.Ingredient) error {
	if ingredient.Name == "" {
		return errors.New("название ингредиента не может быть пустым")
	}
	if ingredient.Quantity <= 0 {
		return errors.New("количество ингредиента должно быть больше нуля")
	}

	return r.db.WithContext(ctx).Create(ingredient).Error
}

// GetByID получает продукт по ID
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// Update обновляет продукт
func (r *ProductRepository) Update(ctx context.Context, product *entity.Product) error {
	if product.ID == uuid.Nil {
		return errors.New("ID продукта не может быть пустым")
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

	return r.db.WithContext(ctx).Save(product).Error
}

// Delete удаляет продукт по ID
func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Product{}).Error
}

// GetByCategory получает продукты по категории
func (r *ProductRepository) GetByCategory(ctx context.Context, category string) ([]*entity.Product, error) {
	var products []*entity.Product
	if err := r.db.WithContext(ctx).Where("category = ?", category).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetActive получает только активные продукты
func (r *ProductRepository) GetActive(ctx context.Context) ([]*entity.Product, error) {
	var products []*entity.Product
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetNotActive получает неактивные продукты
func (r *ProductRepository) GetNotActive(ctx context.Context) ([]*entity.Product, error) {
	var products []*entity.Product
	if err := r.db.WithContext(ctx).Where("is_active = ?", false).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// Count подсчитывает общее количество продуктов
func (r *ProductRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Product{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetAll получает все продукты
func (r *ProductRepository) GetAll(ctx context.Context) ([]*entity.Product, error) {
	var products []*entity.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// Search ищет продукты по названию или описанию
func (r *ProductRepository) Search(ctx context.Context, query string) ([]*entity.Product, error) {
	var products []*entity.Product
	if err := r.db.WithContext(ctx).Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetByPriceRange получает продукты в диапазоне цен
func (r *ProductRepository) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*entity.Product, error) {
	var products []*entity.Product
	if err := r.db.WithContext(ctx).Where("price BETWEEN ? AND ?", minPrice, maxPrice).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// Activate активирует продукт
func (r *ProductRepository) Activate(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	return r.db.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", id).Update("is_active", true).Error
}

// Deactivate деактивирует продукт
func (r *ProductRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	return r.db.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", id).Update("is_active", false).Error
}

// Exists проверяет существование продукта по ID
func (r *ProductRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetIngredientsByProduct получает ингредиенты продукта
func (r *ProductRepository) GetIngredientsByProduct(ctx context.Context, productID uuid.UUID) ([]*entity.Ingredient, error) {
	var ingredients []*entity.Ingredient
	if err := r.db.WithContext(ctx).
		Joins("JOIN product_ingredients ON ingredients.id = product_ingredients.ingredient_id").
		Where("product_ingredients.product_id = ?", productID).
		Find(&ingredients).Error; err != nil {
		return nil, err
	}
	return ingredients, nil
}

// AddIngredientToProduct добавляет ингредиент к продукту
func (r *ProductRepository) AddIngredientToProduct(ctx context.Context, productID, ingredientID uuid.UUID, quantity float64, unit string) error {
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

	// Проверяем, существует ли продукт
	_, err := r.GetByID(ctx, productID)
	if err != nil {
		return errors.New("продукт не найден")
	}

	// Проверяем, существует ли ингредиент
	var ingredient entity.Ingredient
	if err := r.db.WithContext(ctx).Where("id = ?", ingredientID).First(&ingredient).Error; err != nil {
		return errors.New("ингредиент не найден")
	}

	// Проверяем, не добавлен ли уже этот ингредиент к продукту
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.ProductIngredient{}).
		Where("product_id = ? AND ingredient_id = ?", productID, ingredientID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("ингредиент уже добавлен к продукту")
	}

	// Добавляем связь
	productIngredient := &entity.ProductIngredient{
		ProductID:    productID,
		IngredientID: ingredientID,
		Quantity:     quantity,
		Unit:         unit,
	}

	return r.db.WithContext(ctx).Create(productIngredient).Error
}

// RemoveIngredientFromProduct удаляет ингредиент из продукта
func (r *ProductRepository) RemoveIngredientFromProduct(ctx context.Context, productID, ingredientID uuid.UUID) error {
	if productID == uuid.Nil {
		return errors.New("ID продукта не может быть пустым")
	}
	if ingredientID == uuid.Nil {
		return errors.New("ID ингредиента не может быть пустым")
	}

	return r.db.WithContext(ctx).
		Where("product_id = ? AND ingredient_id = ?", productID, ingredientID).
		Delete(&entity.ProductIngredient{}).Error
}

// UpdateProductIngredient обновляет количество ингредиента в продукте
func (r *ProductRepository) UpdateProductIngredient(ctx context.Context, productID, ingredientID uuid.UUID, quantity float64, unit string) error {
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

	return r.db.WithContext(ctx).
		Model(&entity.ProductIngredient{}).
		Where("product_id = ? AND ingredient_id = ?", productID, ingredientID).
		Updates(map[string]interface{}{
			"quantity": quantity,
			"unit":     unit,
		}).Error
}

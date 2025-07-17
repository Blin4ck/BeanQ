package repositories

import (
	"coffe/internal/menu/entity"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IngredientRepository реализует методы доступа к ингредиентам в базе данных.
type IngredientRepository struct {
	db *gorm.DB
}

// NewIngredientRepository создает новый экземпляр IngredientRepository.
func NewIngredientRepository(db *gorm.DB) *IngredientRepository {
	return &IngredientRepository{db: db}
}

// Create добавляет ингредиент в базу данных.
func (r *IngredientRepository) Create(ctx context.Context, ingredient *entity.Ingredient) error {
	if ingredient.Name == "" {
		return errors.New("название ингредиента не может быть пустым")
	}
	if ingredient.Quantity <= 0 {
		return errors.New("количество ингредиента должно быть больше нуля")
	}
	if ingredient.Unit == "" {
		return errors.New("единица измерения не может быть пустой")
	}

	return r.db.WithContext(ctx).Create(ingredient).Error
}

// GetByName возвращает ингредиент по названию.
func (r *IngredientRepository) GetByName(ctx context.Context, name string) (*entity.Ingredient, error) {
	var ingredient entity.Ingredient
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&ingredient).Error; err != nil {
		return nil, err
	}
	return &ingredient, nil
}

// GetByID возвращает ингредиент по идентификатору.
func (r *IngredientRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Ingredient, error) {
	var ingredient entity.Ingredient
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&ingredient).Error; err != nil {
		return nil, err
	}
	return &ingredient, nil
}

// Update обновляет ингредиент
func (r *IngredientRepository) Update(ctx context.Context, ingredient *entity.Ingredient) error {
	if ingredient.ID == uuid.Nil {
		return errors.New("ID ингредиента не может быть пустым")
	}
	if ingredient.Name == "" {
		return errors.New("название ингредиента не может быть пустым")
	}
	if ingredient.Quantity <= 0 {
		return errors.New("количество ингредиента должно быть больше нуля")
	}

	return r.db.WithContext(ctx).Save(ingredient).Error
}

// Delete удаляет ингредиент по ID
func (r *IngredientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Ingredient{}).Error
}

// GetAll получает все ингредиенты
func (r *IngredientRepository) GetAll(ctx context.Context) ([]*entity.Ingredient, error) {
	var ingredients []*entity.Ingredient
	if err := r.db.WithContext(ctx).Find(&ingredients).Error; err != nil {
		return nil, err
	}
	return ingredients, nil
}

// Search ищет ингредиенты по названию
func (r *IngredientRepository) Search(ctx context.Context, query string) ([]*entity.Ingredient, error) {
	var ingredients []*entity.Ingredient
	if err := r.db.WithContext(ctx).Where("name ILIKE ?", "%"+query+"%").Find(&ingredients).Error; err != nil {
		return nil, err
	}
	return ingredients, nil
}

// GetByUnit получает ингредиенты по единице измерения
func (r *IngredientRepository) GetByUnit(ctx context.Context, unit string) ([]*entity.Ingredient, error) {
	var ingredients []*entity.Ingredient
	if err := r.db.WithContext(ctx).Where("unit = ?", unit).Find(&ingredients).Error; err != nil {
		return nil, err
	}
	return ingredients, nil
}

// Count подсчитывает общее количество ингредиентов
func (r *IngredientRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Ingredient{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Exists проверяет существование ингредиента по ID
func (r *IngredientRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Ingredient{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetProductsByIngredient получает продукты, содержащие определенный ингредиент
func (r *IngredientRepository) GetProductsByIngredient(ctx context.Context, ingredientID uuid.UUID) ([]*entity.Product, error) {
	var products []*entity.Product
	if err := r.db.WithContext(ctx).
		Joins("JOIN product_ingredients ON products.id = product_ingredients.product_id").
		Where("product_ingredients.ingredient_id = ?", ingredientID).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

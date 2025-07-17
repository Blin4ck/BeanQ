package usecase

import (
	"coffe/internal/menu/entity"
	"coffe/internal/menu/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type IngredientUsecase struct {
	ingredientRepo repository.IngredientRepository
}

func NewIngredientUsecase(ingredientRepo repository.IngredientRepository) *IngredientUsecase {
	return &IngredientUsecase{ingredientRepo: ingredientRepo}
}

// Create создает новый ингредиент
func (uc *IngredientUsecase) Create(ctx context.Context, ingredient *entity.Ingredient) error {
	if err := uc.validateIngredient(ingredient); err != nil {
		return err
	}

	// Проверяем, не существует ли уже ингредиент с таким названием
	existing, err := uc.ingredientRepo.GetByName(ctx, ingredient.Name)
	if err == nil && existing != nil {
		return errors.New("ингредиент с таким названием уже существует")
	}

	return uc.ingredientRepo.Create(ctx, ingredient)
}

// GetByID получает ингредиент по ID
func (uc *IngredientUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Ingredient, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID не может быть пустым")
	}

	ingredient, err := uc.ingredientRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("ингредиент не найден")
	}

	return ingredient, nil
}

// GetByName получает ингредиент по названию
func (uc *IngredientUsecase) GetByName(ctx context.Context, name string) (*entity.Ingredient, error) {
	if name == "" {
		return nil, errors.New("название не может быть пустым")
	}

	ingredient, err := uc.ingredientRepo.GetByName(ctx, name)
	if err != nil {
		return nil, errors.New("ингредиент не найден")
	}

	return ingredient, nil
}

// Update обновляет ингредиент
func (uc *IngredientUsecase) Update(ctx context.Context, ingredient *entity.Ingredient) error {
	if err := uc.validateIngredient(ingredient); err != nil {
		return err
	}

	// Проверяем, существует ли ингредиент
	existing, err := uc.ingredientRepo.GetByID(ctx, ingredient.ID)
	if err != nil {
		return errors.New("ингредиент не найден")
	}

	// Проверяем, не конфликтует ли новое название с существующим
	if existing.Name != ingredient.Name {
		duplicate, err := uc.ingredientRepo.GetByName(ctx, ingredient.Name)
		if err == nil && duplicate != nil {
			return errors.New("ингредиент с таким названием уже существует")
		}
	}

	return uc.ingredientRepo.Update(ctx, ingredient)
}

// Delete удаляет ингредиент по ID
func (uc *IngredientUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}

	// Проверяем, существует ли ингредиент
	_, err := uc.ingredientRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("ингредиент не найден")
	}

	return uc.ingredientRepo.Delete(ctx, id)
}

// GetAll получает все ингредиенты
func (uc *IngredientUsecase) GetAll(ctx context.Context) ([]*entity.Ingredient, error) {
	ingredients, err := uc.ingredientRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.New("ошибка при получении ингредиентов")
	}

	return ingredients, nil
}

// Search ищет ингредиенты по названию
func (uc *IngredientUsecase) Search(ctx context.Context, query string) ([]*entity.Ingredient, error) {
	if query == "" {
		return nil, errors.New("поисковый запрос не может быть пустым")
	}

	ingredients, err := uc.ingredientRepo.Search(ctx, query)
	if err != nil {
		return nil, errors.New("ошибка при поиске ингредиентов")
	}

	return ingredients, nil
}

// GetByUnit получает ингредиенты по единице измерения
func (uc *IngredientUsecase) GetByUnit(ctx context.Context, unit string) ([]*entity.Ingredient, error) {
	if unit == "" {
		return nil, errors.New("единица измерения не может быть пустой")
	}

	ingredients, err := uc.ingredientRepo.GetByUnit(ctx, unit)
	if err != nil {
		return nil, errors.New("ошибка при получении ингредиентов по единице измерения")
	}

	return ingredients, nil
}

// Count получает общее количество ингредиентов
func (uc *IngredientUsecase) Count(ctx context.Context) (int64, error) {
	count, err := uc.ingredientRepo.Count(ctx)
	if err != nil {
		return 0, errors.New("ошибка при подсчете ингредиентов")
	}

	return count, nil
}

// Exists проверяет существование ингредиента
func (uc *IngredientUsecase) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, errors.New("ID не может быть пустым")
	}

	exists, err := uc.ingredientRepo.Exists(ctx, id)
	if err != nil {
		return false, errors.New("ошибка при проверке существования ингредиента")
	}

	return exists, nil
}

// GetIngredientProducts получает продукты, содержащие определенный ингредиент
func (uc *IngredientUsecase) GetIngredientProducts(ctx context.Context, ingredientID uuid.UUID) ([]*entity.Product, error) {
	if ingredientID == uuid.Nil {
		return nil, errors.New("ID ингредиента не может быть пустым")
	}

	// Проверяем, существует ли ингредиент
	_, err := uc.ingredientRepo.GetByID(ctx, ingredientID)
	if err != nil {
		return nil, errors.New("ингредиент не найден")
	}

	products, err := uc.ingredientRepo.GetProductsByIngredient(ctx, ingredientID)
	if err != nil {
		return nil, errors.New("ошибка при получении продуктов ингредиента")
	}

	return products, nil
}

// validateIngredient валидирует ингредиент
func (uc *IngredientUsecase) validateIngredient(ingredient *entity.Ingredient) error {
	if ingredient == nil {
		return errors.New("ингредиент не может быть пустым")
	}
	if ingredient.Name == "" {
		return errors.New("название ингредиента не может быть пустым")
	}
	if ingredient.Quantity <= 0 {
		return errors.New("количество ингредиента должно быть больше нуля")
	}
	if ingredient.Unit == "" {
		return errors.New("единица измерения не может быть пустой")
	}

	return nil
}

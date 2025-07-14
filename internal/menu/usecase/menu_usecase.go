package usecase

import (
	"coffe/internal/menu/entity"
	"coffe/internal/menu/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type MenuUsecase struct {
	menuRepo repository.MenuRepository
}

func NewMenuUsecase(menuRepo repository.MenuRepository) *MenuUsecase {
	return &MenuUsecase{menuRepo: menuRepo}
}

// Create создает новое меню
func (u *MenuUsecase) Create(ctx context.Context, menu *entity.Menu) error {
	if err := u.validateMenu(menu); err != nil {
		return err
	}

	return u.menuRepo.Create(ctx, menu)
}

// GetByID получает меню по ID
func (u *MenuUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Menu, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID не может быть пустым")
	}

	menu, err := u.menuRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("меню не найдено")
	}

	return menu, nil
}

// GetActive получает активные меню
func (u *MenuUsecase) GetActive(ctx context.Context) ([]*entity.Menu, error) {
	menus, err := u.menuRepo.GetActive(ctx)
	if err != nil {
		return nil, errors.New("ошибка при получении активных меню")
	}
	return menus, nil
}

// GetAll получает все меню
func (u *MenuUsecase) GetAll(ctx context.Context) ([]*entity.Menu, error) {
	menus, err := u.menuRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.New("ошибка при получении меню")
	}
	return menus, nil
}

// Update обновляет меню
func (u *MenuUsecase) Update(ctx context.Context, menu *entity.Menu) error {
	if err := u.validateMenu(menu); err != nil {
		return err
	}

	// Проверяем, существует ли меню
	_, err := u.menuRepo.GetByID(ctx, menu.ID)
	if err != nil {
		return errors.New("меню не найдено")
	}

	return u.menuRepo.Update(ctx, menu)
}

// Delete удаляет меню по ID
func (u *MenuUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}

	// Проверяем, существует ли меню
	_, err := u.menuRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("меню не найдено")
	}

	return u.menuRepo.Delete(ctx, id)
}

// Activate активирует меню
func (u *MenuUsecase) Activate(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}

	// Проверяем, существует ли меню
	_, err := u.menuRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("меню не найдено")
	}

	return u.menuRepo.Activate(ctx, id)
}

// Deactivate деактивирует меню
func (u *MenuUsecase) Deactivate(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}

	// Проверяем, существует ли меню
	_, err := u.menuRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("меню не найдено")
	}

	return u.menuRepo.Deactivate(ctx, id)
}

// ===== МЕТОДЫ ДЛЯ РАБОТЫ С КАТЕГОРИЯМИ =====

// CreateCategory создает новую категорию
func (u *MenuUsecase) CreateCategory(ctx context.Context, category *entity.MenuCategory) error {
	if err := u.validateCategory(category); err != nil {
		return err
	}

	return u.menuRepo.CreateCategory(ctx, category)
}

// GetCategories получает все категории
func (u *MenuUsecase) GetCategories(ctx context.Context) ([]*entity.MenuCategory, error) {
	categories, err := u.menuRepo.GetCategories(ctx)
	if err != nil {
		return nil, errors.New("ошибка при получении категорий")
	}
	return categories, nil
}

// GetCategoryByID получает категорию по ID
func (u *MenuUsecase) GetCategoryByID(ctx context.Context, id uuid.UUID) (*entity.MenuCategory, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID не может быть пустым")
	}

	category, err := u.menuRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, errors.New("категория не найдена")
	}

	return category, nil
}

// UpdateCategory обновляет категорию
func (u *MenuUsecase) UpdateCategory(ctx context.Context, category *entity.MenuCategory) error {
	if err := u.validateCategory(category); err != nil {
		return err
	}

	// Проверяем, существует ли категория
	_, err := u.menuRepo.GetCategoryByID(ctx, category.ID)
	if err != nil {
		return errors.New("категория не найдена")
	}

	return u.menuRepo.UpdateCategory(ctx, category)
}

// DeleteCategory удаляет категорию по ID
func (u *MenuUsecase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}

	// Проверяем, существует ли категория
	_, err := u.menuRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return errors.New("категория не найдена")
	}

	return u.menuRepo.DeleteCategory(ctx, id)
}

// GetCategoriesByMenu получает категории конкретного меню
func (u *MenuUsecase) GetCategoriesByMenu(ctx context.Context, menuID uuid.UUID) ([]*entity.MenuCategory, error) {
	if menuID == uuid.Nil {
		return nil, errors.New("ID меню не может быть пустым")
	}

	// Проверяем, существует ли меню
	_, err := u.menuRepo.GetByID(ctx, menuID)
	if err != nil {
		return nil, errors.New("меню не найдено")
	}

	categories, err := u.menuRepo.GetCategoriesByMenu(ctx, menuID)
	if err != nil {
		return nil, errors.New("ошибка при получении категорий меню")
	}

	return categories, nil
}

// ===== МЕТОДЫ ДЛЯ РАБОТЫ С ПОЗИЦИЯМИ МЕНЮ =====

// CreateMenuItem создает новую позицию меню
func (u *MenuUsecase) CreateMenuItem(ctx context.Context, item *entity.MenuItem) error {
	if err := u.validateMenuItem(item); err != nil {
		return err
	}

	return u.menuRepo.CreateMenuItem(ctx, item)
}

// GetItemsByCategory получает позиции по категории
func (u *MenuUsecase) GetItemsByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.MenuItem, error) {
	if categoryID == uuid.Nil {
		return nil, errors.New("ID категории не может быть пустым")
	}

	// Проверяем, существует ли категория
	_, err := u.menuRepo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, errors.New("категория не найдена")
	}

	items, err := u.menuRepo.GetItemsByCategory(ctx, categoryID)
	if err != nil {
		return nil, errors.New("ошибка при получении позиций категории")
	}

	return items, nil
}

// GetActiveItems получает активные позиции меню
func (u *MenuUsecase) GetActiveItems(ctx context.Context) ([]*entity.MenuItem, error) {
	items, err := u.menuRepo.GetActiveItems(ctx)
	if err != nil {
		return nil, errors.New("ошибка при получении активных позиций")
	}
	return items, nil
}

// GetMenuItemByID получает позицию меню по ID
func (u *MenuUsecase) GetMenuItemByID(ctx context.Context, id uuid.UUID) (*entity.MenuItem, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID не может быть пустым")
	}

	item, err := u.menuRepo.GetMenuItemByID(ctx, id)
	if err != nil {
		return nil, errors.New("позиция меню не найдена")
	}

	return item, nil
}

// UpdateMenuItem обновляет позицию меню
func (u *MenuUsecase) UpdateMenuItem(ctx context.Context, item *entity.MenuItem) error {
	if err := u.validateMenuItem(item); err != nil {
		return err
	}

	// Проверяем, существует ли позиция
	_, err := u.menuRepo.GetMenuItemByID(ctx, item.ID)
	if err != nil {
		return errors.New("позиция меню не найдена")
	}

	return u.menuRepo.UpdateMenuItem(ctx, item)
}

// DeleteMenuItem удаляет позицию меню по ID
func (u *MenuUsecase) DeleteMenuItem(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}

	// Проверяем, существует ли позиция
	_, err := u.menuRepo.GetMenuItemByID(ctx, id)
	if err != nil {
		return errors.New("позиция меню не найдена")
	}

	return u.menuRepo.DeleteMenuItem(ctx, id)
}

// ActivateMenuItem активирует позицию меню
func (u *MenuUsecase) ActivateMenuItem(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}

	// Проверяем, существует ли позиция
	_, err := u.menuRepo.GetMenuItemByID(ctx, id)
	if err != nil {
		return errors.New("позиция меню не найдена")
	}

	return u.menuRepo.ActivateMenuItem(ctx, id)
}

// DeactivateMenuItem деактивирует позицию меню
func (u *MenuUsecase) DeactivateMenuItem(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}

	// Проверяем, существует ли позиция
	_, err := u.menuRepo.GetMenuItemByID(ctx, id)
	if err != nil {
		return errors.New("позиция меню не найдена")
	}

	return u.menuRepo.DeactivateMenuItem(ctx, id)
}

// GetMenuItemsByMenu получает позиции конкретного меню
func (u *MenuUsecase) GetMenuItemsByMenu(ctx context.Context, menuID uuid.UUID) ([]*entity.MenuItem, error) {
	if menuID == uuid.Nil {
		return nil, errors.New("ID меню не может быть пустым")
	}

	// Проверяем, существует ли меню
	_, err := u.menuRepo.GetByID(ctx, menuID)
	if err != nil {
		return nil, errors.New("меню не найдено")
	}

	items, err := u.menuRepo.GetMenuItemsByMenu(ctx, menuID)
	if err != nil {
		return nil, errors.New("ошибка при получении позиций меню")
	}

	return items, nil
}

// ===== УТИЛИТЫ =====

// Count получает количество меню
func (u *MenuUsecase) Count(ctx context.Context) (int64, error) {
	count, err := u.menuRepo.Count(ctx)
	if err != nil {
		return 0, errors.New("ошибка при подсчете меню")
	}
	return count, nil
}

// Exists проверяет существование меню
func (u *MenuUsecase) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if id == uuid.Nil {
		return false, errors.New("ID не может быть пустым")
	}

	exists, err := u.menuRepo.Exists(ctx, id)
	if err != nil {
		return false, errors.New("ошибка при проверке существования меню")
	}

	return exists, nil
}

// ===== ВАЛИДАЦИЯ =====

// validateMenu валидирует меню
func (u *MenuUsecase) validateMenu(menu *entity.Menu) error {
	if menu == nil {
		return errors.New("меню не может быть пустым")
	}
	if menu.Name == "" {
		return errors.New("название меню не может быть пустым")
	}
	if menu.Description == "" {
		return errors.New("описание меню не может быть пустым")
	}

	return nil
}

// validateCategory валидирует категорию
func (u *MenuUsecase) validateCategory(category *entity.MenuCategory) error {
	if category == nil {
		return errors.New("категория не может быть пустой")
	}
	if category.Name == "" {
		return errors.New("название категории не может быть пустым")
	}
	if category.Description == "" {
		return errors.New("описание категории не может быть пустым")
	}

	return nil
}

// validateMenuItem валидирует позицию меню
func (u *MenuUsecase) validateMenuItem(item *entity.MenuItem) error {
	if item == nil {
		return errors.New("позиция меню не может быть пустой")
	}
	if item.MenuID == uuid.Nil {
		return errors.New("ID меню не может быть пустым")
	}
	if item.ProductID == uuid.Nil {
		return errors.New("ID продукта не может быть пустым")
	}
	if item.CategoryID == uuid.Nil {
		return errors.New("ID категории не может быть пустым")
	}

	return nil
}

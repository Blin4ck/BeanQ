package repositories

import (
	"coffe/internal/menu/entity"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// Create создает новое меню
func (r *MenuRepository) Create(ctx context.Context, menu *entity.Menu) error {
	if menu.Name == "" {
		return errors.New("название меню не может быть пустым")
	}
	if menu.Description == "" {
		return errors.New("описание меню не может быть пустым")
	}

	return r.db.WithContext(ctx).Create(menu).Error
}

// GetByID получает меню по ID
func (r *MenuRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Menu, error) {
	var menu entity.Menu
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&menu).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetActive получает активные меню
func (r *MenuRepository) GetActive(ctx context.Context) ([]*entity.Menu, error) {
	var menus []*entity.Menu
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// GetAll получает все меню
func (r *MenuRepository) GetAll(ctx context.Context) ([]*entity.Menu, error) {
	var menus []*entity.Menu
	if err := r.db.WithContext(ctx).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// Update обновляет меню
func (r *MenuRepository) Update(ctx context.Context, menu *entity.Menu) error {
	if menu.ID == uuid.Nil {
		return errors.New("ID меню не может быть пустым")
	}
	if menu.Name == "" {
		return errors.New("название меню не может быть пустым")
	}

	return r.db.WithContext(ctx).Save(menu).Error
}

// Delete удаляет меню по ID
func (r *MenuRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Menu{}).Error
}

// Activate активирует меню
func (r *MenuRepository) Activate(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	return r.db.WithContext(ctx).Model(&entity.Menu{}).Where("id = ?", id).Update("is_active", true).Error
}

// Deactivate деактивирует меню
func (r *MenuRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	return r.db.WithContext(ctx).Model(&entity.Menu{}).Where("id = ?", id).Update("is_active", false).Error
}

// ===== МЕТОДЫ ДЛЯ РАБОТЫ С КАТЕГОРИЯМИ =====

// CreateCategory создает новую категорию
func (r *MenuRepository) CreateCategory(ctx context.Context, category *entity.MenuCategory) error {
	if category.Name == "" {
		return errors.New("название категории не может быть пустым")
	}
	if category.Description == "" {
		return errors.New("описание категории не может быть пустым")
	}

	return r.db.WithContext(ctx).Create(category).Error
}

// GetCategories получает все категории
func (r *MenuRepository) GetCategories(ctx context.Context) ([]*entity.MenuCategory, error) {
	var categories []*entity.MenuCategory
	if err := r.db.WithContext(ctx).Order("sort_order").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID получает категорию по ID
func (r *MenuRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*entity.MenuCategory, error) {
	var category entity.MenuCategory
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory обновляет категорию
func (r *MenuRepository) UpdateCategory(ctx context.Context, category *entity.MenuCategory) error {
	if category.ID == uuid.Nil {
		return errors.New("ID категории не может быть пустым")
	}
	if category.Name == "" {
		return errors.New("название категории не может быть пустым")
	}

	return r.db.WithContext(ctx).Save(category).Error
}

// DeleteCategory удаляет категорию по ID
func (r *MenuRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.MenuCategory{}).Error
}

// GetCategoriesByMenu получает категории конкретного меню
func (r *MenuRepository) GetCategoriesByMenu(ctx context.Context, menuID uuid.UUID) ([]*entity.MenuCategory, error) {
	var categories []*entity.MenuCategory
	if err := r.db.WithContext(ctx).
		Joins("JOIN menu_categories ON menu_categories.category_id = menu_categories.id").
		Where("menu_categories.menu_id = ?", menuID).
		Order("sort_order").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// ===== МЕТОДЫ ДЛЯ РАБОТЫ С ПОЗИЦИЯМИ МЕНЮ =====

// CreateMenuItem создает новую позицию меню
func (r *MenuRepository) CreateMenuItem(ctx context.Context, item *entity.MenuItem) error {
	if item.MenuID == uuid.Nil {
		return errors.New("ID меню не может быть пустым")
	}
	if item.ProductID == uuid.Nil {
		return errors.New("ID продукта не может быть пустым")
	}
	if item.CategoryID == uuid.Nil {
		return errors.New("ID категории не может быть пустым")
	}

	return r.db.WithContext(ctx).Create(item).Error
}

// GetItemsByCategory получает позиции по категории
func (r *MenuRepository) GetItemsByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.MenuItem, error) {
	var items []*entity.MenuItem
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Category").
		Where("category_id = ?", categoryID).
		Order("sort_order").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetActiveItems получает активные позиции меню
func (r *MenuRepository) GetActiveItems(ctx context.Context) ([]*entity.MenuItem, error) {
	var items []*entity.MenuItem
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Category").
		Where("is_active = ?", true).
		Order("sort_order").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetMenuItemByID получает позицию меню по ID
func (r *MenuRepository) GetMenuItemByID(ctx context.Context, id uuid.UUID) (*entity.MenuItem, error) {
	var item entity.MenuItem
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Category").
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateMenuItem обновляет позицию меню
func (r *MenuRepository) UpdateMenuItem(ctx context.Context, item *entity.MenuItem) error {
	if item.ID == uuid.Nil {
		return errors.New("ID позиции не может быть пустым")
	}

	return r.db.WithContext(ctx).Save(item).Error
}

// DeleteMenuItem удаляет позицию меню по ID
func (r *MenuRepository) DeleteMenuItem(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.MenuItem{}).Error
}

// ActivateMenuItem активирует позицию меню
func (r *MenuRepository) ActivateMenuItem(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	return r.db.WithContext(ctx).Model(&entity.MenuItem{}).Where("id = ?", id).Update("is_active", true).Error
}

// DeactivateMenuItem деактивирует позицию меню
func (r *MenuRepository) DeactivateMenuItem(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	return r.db.WithContext(ctx).Model(&entity.MenuItem{}).Where("id = ?", id).Update("is_active", false).Error
}

// GetMenuItemsByMenu получает позиции конкретного меню
func (r *MenuRepository) GetMenuItemsByMenu(ctx context.Context, menuID uuid.UUID) ([]*entity.MenuItem, error) {
	var items []*entity.MenuItem
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Category").
		Where("menu_id = ?", menuID).
		Order("sort_order").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ===== УТИЛИТЫ =====

// Count подсчитывает общее количество меню
func (r *MenuRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Menu{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Exists проверяет существование меню по ID
func (r *MenuRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Menu{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

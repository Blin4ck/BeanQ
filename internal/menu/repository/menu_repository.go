package repository

import (
	"coffe/internal/menu/entity"
	"context"

	"github.com/google/uuid"
)

type MenuRepository interface {
	// Методы для работы с меню
	Create(ctx context.Context, menu *entity.Menu) error             // создание меню
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Menu, error) // меню по id
	GetActive(ctx context.Context) ([]*entity.Menu, error)           // активные меню
	GetAll(ctx context.Context) ([]*entity.Menu, error)              // все меню
	Update(ctx context.Context, menu *entity.Menu) error             // обновление меню
	Delete(ctx context.Context, id uuid.UUID) error                  // удаление меню
	Activate(ctx context.Context, id uuid.UUID) error                // активировать меню
	Deactivate(ctx context.Context, id uuid.UUID) error              // деактивировать меню

	// Методы для работы с категориями
	CreateCategory(ctx context.Context, category *entity.MenuCategory) error                   // создание категории
	GetCategories(ctx context.Context) ([]*entity.MenuCategory, error)                         // все категории
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*entity.MenuCategory, error)           // категория по id
	UpdateCategory(ctx context.Context, category *entity.MenuCategory) error                   // обновление категории
	DeleteCategory(ctx context.Context, id uuid.UUID) error                                    // удаление категории
	GetCategoriesByMenu(ctx context.Context, menuID uuid.UUID) ([]*entity.MenuCategory, error) // категории меню

	// Методы для работы с позициями меню
	CreateMenuItem(ctx context.Context, item *entity.MenuItem) error                          // создание позиции меню
	GetItemsByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.MenuItem, error) // товары в категории
	GetActiveItems(ctx context.Context) ([]*entity.MenuItem, error)                           // активные позиции меню
	GetMenuItemByID(ctx context.Context, id uuid.UUID) (*entity.MenuItem, error)              // позиция по id
	UpdateMenuItem(ctx context.Context, item *entity.MenuItem) error                          // обновление позиции
	DeleteMenuItem(ctx context.Context, id uuid.UUID) error                                   // удаление позиции
	ActivateMenuItem(ctx context.Context, id uuid.UUID) error                                 // активировать позицию
	DeactivateMenuItem(ctx context.Context, id uuid.UUID) error                               // деактивировать позицию
	GetMenuItemsByMenu(ctx context.Context, menuID uuid.UUID) ([]*entity.MenuItem, error)     // позиции меню

	// Утилиты
	Count(ctx context.Context) (int64, error)               // количество меню
	Exists(ctx context.Context, id uuid.UUID) (bool, error) // проверить существование
}

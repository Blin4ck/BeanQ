package repository

import (
	"coffe/internal/order/entity"
	"context"

	"github.com/google/uuid"
)

// OrderRepository определяет методы для работы с заказами.
type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error                                // создание заказа
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)                     // поиск по id
	Update(ctx context.Context, order *entity.Order) error                                // обновление заказа
	Delete(ctx context.Context, id uuid.UUID) error                                       // удаление заказа
	GetByCustomer(ctx context.Context, customerID uuid.UUID) ([]*entity.Order, error)     // заказы клиента
	GetByStatus(ctx context.Context, status entity.OrderStatus) ([]*entity.Order, error)  // заказы по статусу
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status entity.OrderStatus) error // обновление статуса
	Count(ctx context.Context) (int64, error)                                             // количество заказов
	GetToday(ctx context.Context) ([]*entity.Order, error)                                // заказы за сегодня
}

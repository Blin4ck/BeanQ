package usecase

import (
	"coffe/internal/order/entity"
	"coffe/internal/order/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

// OrderUsecase реализует бизнес-логику для работы с заказами.
type OrderUsecase struct {
	orderRepo repository.OrderRepository
}

// NewOrderUsecase создает новый экземпляр OrderUsecase.
func NewOrderUsecase(orderRepo repository.OrderRepository) *OrderUsecase {
	return &OrderUsecase{orderRepo: orderRepo}
}

// Create создает новый заказ.
func (u *OrderUsecase) Create(ctx context.Context, order *entity.Order) error {
	if order.CustomerID == uuid.Nil {
		return errors.New("customer_id не может быть пустым")
	}
	if len(order.Items) == 0 {
		return errors.New("заказ не может быть пустым")
	}
	return u.orderRepo.Create(ctx, order)
}

// GetByID возвращает заказ по идентификатору.
func (u *OrderUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	if id == uuid.Nil {
		return nil, errors.New("id не может быть пустым")
	}
	return u.orderRepo.GetByID(ctx, id)
}

// Update обновляет заказ.
func (u *OrderUsecase) Update(ctx context.Context, order *entity.Order) error {
	if order.Id == uuid.Nil {
		return errors.New("id не может быть пустым")
	}
	return u.orderRepo.Update(ctx, order)
}

// Delete удаляет заказ по идентификатору.
func (u *OrderUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id не может быть пустым")
	}
	return u.orderRepo.Delete(ctx, id)
}

// GetByCustomer возвращает заказы клиента.
func (u *OrderUsecase) GetByCustomer(ctx context.Context, customerID uuid.UUID) ([]*entity.Order, error) {
	if customerID == uuid.Nil {
		return nil, errors.New("customer_id не может быть пустым")
	}
	return u.orderRepo.GetByCustomer(ctx, customerID)
}

// GetByStatus возвращает заказы по статусу.
func (u *OrderUsecase) GetByStatus(ctx context.Context, status entity.OrderStatus) ([]*entity.Order, error) {
	return u.orderRepo.GetByStatus(ctx, status)
}

func (u *OrderUsecase) UpdateStatus(ctx context.Context, orderID uuid.UUID, status entity.OrderStatus) error {
	if orderID == uuid.Nil {
		return errors.New("order_id не может быть пустым")
	}
	return u.orderRepo.UpdateStatus(ctx, orderID, status)
}

func (u *OrderUsecase) Count(ctx context.Context) (int64, error) {
	return u.orderRepo.Count(ctx)
}

func (u *OrderUsecase) GetToday(ctx context.Context) ([]*entity.Order, error) {
	return u.orderRepo.GetToday(ctx)
}

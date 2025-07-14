package repositories

import (
	"coffe/internal/order/entity"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *entity.Order) error {
	if order.CustomerID == uuid.Nil {
		return errors.New("customer_ID не может быть пустым")
	}
	if len(order.Items) == 0 {
		return errors.New("заказ не может быть пустым")
	}
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var order *entity.Order
	if id == uuid.Nil {
		return nil, errors.New("передан пустой id")
	}
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, errors.New("заказ с таким id не найден")
	}
	return order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *entity.Order) error {
	if order.Id == uuid.Nil {
		return errors.New("id не может быть пустым")
	}
	return r.db.WithContext(ctx).Save(order).Error
}

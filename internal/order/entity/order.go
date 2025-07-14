package entity

import (
	"coffe/internal/common"
	"time"

	"github.com/google/uuid"
)

// PaymentMethod определяет способ оплаты заказа.
type PaymentMethod string

const (
	PaymentMethodCash   PaymentMethod = "Наличка"
	PaymentMethodCard   PaymentMethod = "Карта"
	PaymentMethodOnline PaymentMethod = "Иное"
)

// OrderStatus определяет статус заказа.
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "ожидает"
	OrderStatusConfirmed OrderStatus = "подтвержден"
	OrderStatusPreparing OrderStatus = "готовится"
	OrderStatusReady     OrderStatus = "готов"
	OrderStatusCompleted OrderStatus = "выполнен"
	OrderStatusCancelled OrderStatus = "отменен"
)

// Order представляет заказ клиента.
type Order struct {
	Id            uuid.UUID     `json:"id" db:"id"`
	CustomerID    uuid.UUID     `json:"customer_id" db:"customer_id"`
	Customer      *common.User  `json:"customer,omitempty" db:"customer"`
	Items         []ItemsOrders `json:"items" db:"items"`
	Status        OrderStatus   `json:"status" db:"status"`
	Notes         string        `json:"notes" db:"notes"`
	TotalPrice    float64       `json:"total_price" db:"total_price"`
	PaymentMethod PaymentMethod `json:"payment_method" db:"payment_method"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// ItemsOrders представляет позицию заказа.
type ItemsOrders struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	OrderID   uuid.UUID       `json:"order_id" db:"order_id"`
	ProductID uuid.UUID       `json:"product_id" db:"product_id"`
	Product   *common.Product `json:"product,omitempty" db:"product"`
	Quantity  int             `json:"quantity" db:"quantity"`
	Price     float64         `json:"price" db:"price"` // цена на момент заказа
}

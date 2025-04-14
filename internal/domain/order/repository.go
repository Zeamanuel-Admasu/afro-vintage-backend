package order

import (
	"context"
)

type Repository interface {
	CreateOrder(ctx context.Context, o *Order) error
	GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status OrderStatus) error
	DeleteOrder(ctx context.Context, orderID string) error
}

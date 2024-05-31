package api

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/andranikuz/gophermart/pkg/domain/order"
)

type OrderServiceInterface interface {
	SetOrder(ctx context.Context, id uuid.UUID, orderNumber int, userID *uuid.UUID) error
	UserOrders(ctx context.Context, userID *uuid.UUID) ([]*order.Order, error)
	UpdateOrderStatus(ctx context.Context, number int, status order.OrderStatus) error
}

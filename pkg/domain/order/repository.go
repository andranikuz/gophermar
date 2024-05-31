package order

import (
	"context"

	"github.com/gofrs/uuid"
)

type Repository interface {
	Insert(ctx context.Context, order *Order) error
	UpdateOrderStatus(ctx context.Context, number int, status OrderStatus) error
	GetByNumber(ctx context.Context, number int) (*Order, error)
	ListByUserID(ctx context.Context, userID *uuid.UUID) ([]*Order, error)
	ListByStatuses(ctx context.Context, statuses []OrderStatus) ([]*Order, error)
}

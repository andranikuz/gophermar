package order

import (
	"context"

	"github.com/andranikuz/gophermart/pkg/domain/order"
)

func (s OrderService) OrdersByStatuses(
	ctx context.Context,
	statuses []order.OrderStatus,
) ([]*order.Order, error) {
	return s.repo.ListByStatuses(ctx, statuses)
}

package order

import (
	"context"
	"github.com/andranikuz/gophermart/pkg/domain/order"
)

func (s OrderService) UpdateOrderStatus(ctx context.Context, number int, status order.OrderStatus) error {
	return s.repo.UpdateOrderStatus(ctx, number, status)
}

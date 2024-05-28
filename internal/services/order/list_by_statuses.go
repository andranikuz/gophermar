package order

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/andranikuz/gophermart/pkg/domain/order"
)

func (s OrderService) UserOrders(
	ctx context.Context,
	userID *uuid.UUID,
) ([]*order.Order, error) {
	return s.repo.ListByUserID(ctx, userID)
}

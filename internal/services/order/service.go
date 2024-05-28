package order

import (
	"github.com/andranikuz/gophermart/pkg/domain/order"
)

type OrderService struct {
	repo order.Repository
}

func NewOrderService(orderRepo order.Repository) *OrderService {
	return &OrderService{
		repo: orderRepo,
	}
}

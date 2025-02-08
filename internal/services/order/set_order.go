package order

import (
	"context"
	"errors"
	"github.com/andranikuz/gophermart/pkg/domain/order"
	"time"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

var ErrAccrualTransactionCreatedByAnotherUser = errors.New("order: new order: already created by another user")

var ErrAccrualTransactionCreatedBySameUser = errors.New("order: new order: already created by same user")

func (s OrderService) SetOrder(ctx context.Context, id uuid.UUID, orderNumber int, userID *uuid.UUID) error {
	o, _ := s.repo.GetByNumber(ctx, orderNumber)
	if o != nil {
		if o.UserID.String() == userID.String() {
			return ErrAccrualTransactionCreatedBySameUser
		} else {
			return ErrAccrualTransactionCreatedByAnotherUser
		}
	}
	o = &order.Order{
		ID:        id,
		Number:    orderNumber,
		UserID:    userID,
		Status:    order.OrderStatusNew,
		CreatedAt: time.Now(),
	}
	err := s.repo.Insert(ctx, o)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return nil
}

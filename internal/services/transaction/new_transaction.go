package transaction

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"

	"github.com/andranikuz/gophermart/pkg/domain/transaction"
)

func (s TransactionService) NewTransaction(
	ctx context.Context,
	id uuid.UUID,
	orderNumber int,
	transactionType transaction.TransactionType,
	userID *uuid.UUID,
	sum float64,
) error {
	t := &transaction.Transaction{
		ID:          id,
		OrderNumber: orderNumber,
		UserID:      userID,
		Type:        transactionType,
		Amount:      sum,
		CreatedAt:   time.Now(),
	}
	err := s.repo.Insert(ctx, t)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return nil
}

package transaction

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/andranikuz/gophermart/pkg/domain/transaction"
)

type UserBalance struct {
	Current   float64
	Withdrawn float64
}

func (s TransactionService) UserBalance(ctx context.Context, userID *uuid.UUID) (*UserBalance, error) {
	accruals, err := s.repo.UserTransactionsByType(ctx, userID, transaction.TransactionTypeAccrual)
	if err != nil {
		return nil, err
	}
	var aSum float64
	for _, t := range accruals {
		aSum += t.Amount
	}

	withdrawals, err := s.repo.UserTransactionsByType(ctx, userID, transaction.TransactionTypeWithdrawal)
	if err != nil {
		return nil, err
	}
	var wSum float64
	for _, t := range withdrawals {
		wSum += t.Amount
	}

	return &UserBalance{
		Current:   aSum - wSum,
		Withdrawn: wSum,
	}, nil
}

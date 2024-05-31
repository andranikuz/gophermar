package api

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/andranikuz/gophermart/pkg/domain/transaction"
)

type UserBalance struct {
	Current   float64
	Withdrawn float64
}

type TransactionServiceInterface interface {
	NewTransaction(ctx context.Context,
		id uuid.UUID,
		orderNumber int,
		transactionType transaction.TransactionType,
		userID *uuid.UUID,
		sum float64,
	) error
	UserBalance(ctx context.Context, userID *uuid.UUID) (*UserBalance, error)
	UserTransactionsByType(
		ctx context.Context,
		userID *uuid.UUID,
		t transaction.TransactionType,
	) ([]transaction.Transaction, error)
}

package transaction

import (
	"context"

	"github.com/gofrs/uuid"
)

type Repository interface {
	GetByOrderNumber(ctx context.Context, orderNumber int) (*Transaction, error)
	Insert(ctx context.Context, user *Transaction) error
	UserTransactionsByType(ctx context.Context, userID *uuid.UUID, t TransactionType) ([]Transaction, error)
}

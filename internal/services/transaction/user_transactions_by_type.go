package transaction

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/andranikuz/gophermart/pkg/domain/transaction"
)

func (s TransactionService) UserTransactionsByType(
	ctx context.Context,
	userID *uuid.UUID,
	t transaction.TransactionType,
) ([]transaction.Transaction, error) {
	return s.repo.UserTransactionsByType(ctx, userID, t)
}

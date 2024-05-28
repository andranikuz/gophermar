package transaction

import (
	"github.com/andranikuz/gophermart/pkg/domain/transaction"
)

type TransactionService struct {
	repo transaction.Repository
}

func NewTransactionService(userRepo transaction.Repository) *TransactionService {
	return &TransactionService{
		repo: userRepo,
	}
}

package transaction

import (
	"time"

	"github.com/gofrs/uuid"
)

type Transaction struct {
	ID          uuid.UUID
	OrderNumber int
	UserID      *uuid.UUID
	Type        TransactionType
	Amount      float64
	CreatedAt   time.Time
}

type TransactionType string

const (
	TransactionTypeAccrual    TransactionType = "accrual"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
)

package postgres

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/andranikuz/gophermart/pkg/domain/transaction"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r TransactionRepository) CreateTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS "transaction" (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			order_number BIGINT,
			type VARCHAR(16),
			amount float,
			created_at TIMESTAMP
		);`
	_, err := r.db.Exec(query)

	return err
}

func (r TransactionRepository) GetByOrderNumber(ctx context.Context, orderNumber int) (*transaction.Transaction, error) {
	t := transaction.Transaction{}
	query := `SELECT id, user_id, order_number, amount, created_at
			FROM public.transaction 
			WHERE order_number = $1 
		`
	err := r.db.QueryRowContext(ctx, query, orderNumber).Scan(
		&t.ID,
		&t.UserID,
		&t.OrderNumber,
		&t.Type,
		&t.Amount,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r TransactionRepository) Insert(ctx context.Context, transaction *transaction.Transaction) error {
	query := `INSERT INTO public.transaction 
    		(id, user_id, order_number, type, amount, created_at) 
			VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := r.db.ExecContext(
		ctx,
		query,
		transaction.ID,
		transaction.UserID,
		transaction.OrderNumber,
		transaction.Type,
		transaction.Amount,
		transaction.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r TransactionRepository) UserTransactionsByType(
	ctx context.Context,
	userID *uuid.UUID,
	t transaction.TransactionType,
) ([]transaction.Transaction, error) {
	var transactions []transaction.Transaction
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, order_number, type, amount, created_at
				FROM public.transaction 
				WHERE user_id = $1 
				AND type = $2
				ORDER BY created_at
		`,
		userID,
		t,
	)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	for rows.Next() {
		var transaction transaction.Transaction
		if err = rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.OrderNumber,
			&transaction.Type,
			&transaction.Amount,
			&transaction.CreatedAt,
		); err != nil {
			return transactions, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

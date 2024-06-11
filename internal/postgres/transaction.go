package postgres

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andranikuz/gophermart/pkg/domain/transaction"
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		pool: pool,
	}
}

func (r TransactionRepository) CreateTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS "transaction" (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			order_number BIGINT,
			type VARCHAR(16),
			amount float,
			created_at TIMESTAMP
		);`
	_, err := r.pool.Exec(ctx, query)

	return err
}

func (r TransactionRepository) GetByOrderNumber(ctx context.Context, orderNumber int) (*transaction.Transaction, error) {
	t := transaction.Transaction{}
	query := `SELECT id, user_id, order_number, amount, created_at
			FROM public.transaction 
			WHERE order_number = $1 
		`
	err := r.pool.QueryRow(ctx, query, orderNumber).Scan(
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

	if _, err := r.pool.Exec(
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
	rows, err := r.pool.Query(
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

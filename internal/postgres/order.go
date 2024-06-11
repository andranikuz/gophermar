package postgres

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andranikuz/gophermart/pkg/domain/order"
	"github.com/andranikuz/gophermart/pkg/domain/transaction"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepositoryRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		pool: pool,
	}
}

func (r OrderRepository) CreateTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS "order" (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			status varchar(15),
			number BIGINT,
			created_at TIMESTAMP
		);`
	_, err := r.pool.Exec(ctx, query)

	return err
}

func (r OrderRepository) GetByNumber(ctx context.Context, number int) (*order.Order, error) {
	o := order.Order{}
	query := `SELECT id, user_id, number, status, created_at
			FROM public.order 
			WHERE number = $1 
		`
	err := r.pool.QueryRow(ctx, query, number).Scan(
		&o.ID,
		&o.UserID,
		&o.Number,
		&o.Status,
		&o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (r OrderRepository) Insert(ctx context.Context, order *order.Order) error {
	query := `INSERT INTO public.order 
    		(id, user_id, number, status, created_at) 
			VALUES ($1, $2, $3, $4, $5)`

	if _, err := r.pool.Exec(
		ctx,
		query,
		order.ID,
		order.UserID,
		order.Number,
		order.Status,
		order.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r OrderRepository) ListByUserID(ctx context.Context, userID *uuid.UUID) ([]*order.Order, error) {
	var orders []*order.Order
	rows, err := r.pool.Query(
		ctx,
		`SELECT o.id, o.user_id, o.number, o.status, coalesce(SUM(t.amount), 0), o.created_at
				FROM public.order o 
				LEFT JOIN transaction t ON t.order_number = o.number AND t.type = $1
				WHERE o.user_id = $2
				GROUP BY o.id
				ORDER BY o.created_at
		`,
		transaction.TransactionTypeAccrual,
		userID,
	)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	for rows.Next() {
		var o order.Order
		if err = rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Number,
			&o.Status,
			&o.Accrual,
			&o.CreatedAt,
		); err != nil {
			return orders, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

func (r OrderRepository) UpdateOrderStatus(ctx context.Context, number int, status order.OrderStatus) error {
	query := `
		UPDATE public.order 
		SET status = $1
		WHERE number = $2
		`

	if _, err := r.pool.Exec(
		ctx,
		query,
		status,
		number,
	); err != nil {
		return err
	}

	return nil
}

func (r OrderRepository) ListByStatuses(ctx context.Context, statuses []order.OrderStatus) ([]*order.Order, error) {
	var orders []*order.Order
	rows, err := r.pool.Query(
		ctx,
		`SELECT id, user_id, number, status, created_at
				FROM public.order
				WHERE status in $1
		`,
		transaction.TransactionTypeWithdrawal,
		statuses,
	)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	for rows.Next() {
		var o order.Order
		if err = rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Number,
			&o.Status,
			&o.Accrual,
			&o.CreatedAt,
		); err != nil {
			return orders, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

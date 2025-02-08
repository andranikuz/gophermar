package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andranikuz/gophermart/pkg/domain/user"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r UserRepository) CreateTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS public."user" (
		   id uuid NOT NULL,
		   login varchar NOT NULL,
		   "password_hash" varchar NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS user_login_idx ON public."user" USING btree (login)
	;`
	_, err := r.pool.Exec(ctx, query)

	return err
}

func (r UserRepository) Get(ctx context.Context, login string) (*user.User, error) {
	c := user.User{}

	query := `SELECT id, login, password_hash FROM public.user WHERE login = $1`

	err := r.pool.QueryRow(ctx, query, login).Scan(&c.ID, &c.Login, &c.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r UserRepository) Insert(ctx context.Context, user *user.User) error {
	query := `INSERT INTO public.user (id, login, password_hash) VALUES ($1, $2, $3)`

	if _, err := r.pool.Exec(ctx, query, user.ID, user.Login, user.PasswordHash); err != nil {
		return err
	}

	return nil
}

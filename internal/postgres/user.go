package postgres

import (
	"context"
	"database/sql"

	"github.com/andranikuz/gophermart/pkg/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) CreateTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS public."user" (
		   id uuid NOT NULL,
		   login varchar NOT NULL,
		   "password" varchar NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS user_login_idx ON public."user" USING btree (login)
	;`
	_, err := r.db.Exec(query)

	return err
}

func (r UserRepository) Get(ctx context.Context, login string) (*user.User, error) {
	c := user.User{}

	query := `SELECT id, login, password FROM public.user WHERE login = $1`

	err := r.db.QueryRowContext(ctx, query, login).Scan(&c.ID, &c.Login, &c.Password)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r UserRepository) Insert(ctx context.Context, user *user.User) error {
	query := `INSERT INTO public.user (id, login, password) VALUES ($1, $2, $3)`

	if _, err := r.db.ExecContext(ctx, query, user.ID, user.Login, user.Password); err != nil {
		return err
	}

	return nil
}

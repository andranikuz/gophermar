package user

import (
	"context"
)

type Repository interface {
	Get(ctx context.Context, login string) (*User, error)
	Insert(ctx context.Context, user *User) error
}

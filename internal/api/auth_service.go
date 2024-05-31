package api

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/andranikuz/gophermart/pkg/domain/user"
)

type AuthenticationServiceInterface interface {
	Login(ctx context.Context, login string, password string) (*user.User, error)
	Register(ctx context.Context, id uuid.UUID, login string, password string) error
	Token(userID uuid.UUID) (string, error)
	ParseToken(tokenString string) (*uuid.UUID, error)
}

package auth

import (
	"github.com/andranikuz/gophermart/pkg/domain/user"
)

type AuthenticationService struct {
	userRepo user.Repository
}

func NewAuthService(userRepo user.Repository) *AuthenticationService {
	return &AuthenticationService{
		userRepo: userRepo,
	}
}

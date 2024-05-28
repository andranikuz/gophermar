package auth

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"

	"github.com/andranikuz/gophermart/pkg/domain/user"
	"github.com/andranikuz/gophermart/pkg/utils"
)

var ErrUserExists = errors.New("auth: user already exists")

func (auth AuthenticationService) Register(ctx context.Context, id uuid.UUID, login string, password string) error {
	u, _ := auth.userRepo.Get(ctx, login)
	if u != nil {
		log.Info().Msg(ErrUserExists.Error())
		return ErrUserExists
	}
	// hash password
	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Error().Msg(`auth register: ` + err.Error())
		return err
	}

	u = &user.User{
		ID:       id,
		Login:    login,
		Password: hash,
	}

	err = auth.userRepo.Insert(ctx, u)
	if err != nil {
		log.Error().Msg(`auth: register: ` + err.Error())
		return err
	}

	return nil
}

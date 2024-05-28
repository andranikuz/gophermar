package auth

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/andranikuz/gophermart/pkg/domain/user"
	"github.com/andranikuz/gophermart/pkg/utils"
)

var ErrWrongCredentials = errors.New("auth: wrong credentials")

func (auth AuthenticationService) Login(ctx context.Context, login string, password string) (*user.User, error) {
	u, err := auth.userRepo.Get(ctx, login)
	if err != nil {
		log.Error().Msg(`auth: login: ` + err.Error())
		return nil, err
	}
	if utils.CheckPasswordHash(password, u.Password) == false {
		log.Info().Msg(ErrWrongCredentials.Error())
		return nil, ErrWrongCredentials
	}

	return u, nil
}

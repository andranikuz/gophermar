package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/andranikuz/gophermart/pkg/domain/user"
	"github.com/andranikuz/gophermart/pkg/utils"
)

var ErrWrongCredentials = errors.New("auth: wrong credentials")

func (auth AuthenticationService) Login(ctx context.Context, login string, password string) (*user.User, error) {
	u, err := auth.userRepo.Get(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info().Msg(err.Error())
			return nil, ErrWrongCredentials
		} else {
			log.Error().Msg(`auth: login: ` + err.Error())
			return nil, err
		}
	}
	if !utils.CheckPasswordHash(password, u.PasswordHash) {
		log.Info().Msg(ErrWrongCredentials.Error())
		return nil, ErrWrongCredentials
	}

	return u, nil
}

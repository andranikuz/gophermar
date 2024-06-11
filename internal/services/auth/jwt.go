package auth

import (
	"errors"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/andranikuz/gophermart/internal/config"
)

type claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

var ErrTokenNotValid = errors.New("jwt: token not valid")

func (auth *AuthenticationService) Token(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256,
		claims{
			UserID: userID,
		})
	tokenString, err := token.SignedString(config.Config.PrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (auth *AuthenticationService) ParseToken(tokenString string) (*uuid.UUID, error) {
	var c claims
	token, err := jwt.ParseWithClaims(tokenString, &c, func(token *jwt.Token) (interface{}, error) {
		return config.Config.PublicKey, nil
	})
	if err != nil {
		log.Info().Msg(err.Error())
		return nil, err
	}
	if !token.Valid {
		log.Info().Msg(ErrTokenNotValid.Error())
		return nil, ErrTokenNotValid
	}

	return &c.UserID, nil
}

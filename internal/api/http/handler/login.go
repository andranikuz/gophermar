package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"

	"github.com/andranikuz/gophermart/internal/services/auth"
)

type loginRequest struct {
	Login    string `json:"login" validate:"max=50,min=1"`
	Password string `json:"password" validate:"max=50,min=1"`
}

func (h HTTPHandler) LoginHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	body, err := io.ReadAll(r.Body)
	defer logErrorIfExists(err)
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check validation
	validationErr := validator.New().Struct(req)
	if validationErr != nil {
		log.Info().Msg(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, loginErr := h.authenticationService.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(loginErr, auth.ErrWrongCredentials) {
			log.Info().Msg(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			logErrorIfExists(loginErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err = h.SetSession(u.ID, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

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
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &req); err != nil {
		log.Info().Msg(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check validation
	err := validator.New().Struct(req)
	if err != nil {
		log.Info().Msg(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := h.authenticationService.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrWrongCredentials) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err = h.SetSession(u.ID, w); err != nil {
		log.Error().Msg(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

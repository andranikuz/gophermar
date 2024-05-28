package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"

	"github.com/andranikuz/gophermart/internal/services/auth"
)

type registerRequest struct {
	Login    string `json:"login" validate:"max=50,min=1"`
	Password string `json:"password" validate:"max=50,min=1"`
}

func (h HTTPHandler) RegisterHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req registerRequest
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
	userID, _ := uuid.NewV6()
	// handle request
	err = h.authenticationService.Register(ctx, userID, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.SetSession(userID, w); err != nil {
		log.Error().Msg(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

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
	body, err := io.ReadAll(r.Body)
	defer logErrorIfExists(err)
	if validationErr := json.Unmarshal(body, &req); err != nil {
		log.Info().Msg(validationErr.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check validation
	validationErr := validator.New().Struct(req)
	if validationErr != nil {
		log.Info().Msg(validationErr.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, _ := uuid.NewV6()
	// handle request
	registerErr := h.authenticationService.Register(ctx, userID, req.Login, req.Password)
	if registerErr != nil {
		if errors.Is(registerErr, auth.ErrUserExists) {
			log.Info().Msg(registerErr.Error())
			w.WriteHeader(http.StatusConflict)
			return
		}
		logErrorIfExists(registerErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.SetSession(userID, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

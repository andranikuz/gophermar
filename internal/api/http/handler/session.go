package handler

import (
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

var ErrUserNotAuthed = errors.New("session: user not authed")

func (h HTTPHandler) SetSession(userID uuid.UUID, w http.ResponseWriter) error {
	token, err := h.authenticationService.Token(userID)
	if err != nil {
		log.Error().Msg(`http token generate: ` + err.Error())
		return err
	}

	cookie := &http.Cookie{Name: "Authorization", Value: token}
	http.SetCookie(w, cookie)
	return nil
}

func (h HTTPHandler) authMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, err := h.GetUserID(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (h HTTPHandler) GetUserID(r *http.Request) (*uuid.UUID, error) {
	cookie, err := r.Cookie("Authorization")
	if err != nil || cookie.Value == "" {
		return nil, ErrUserNotAuthed
	}

	userID, err := h.authenticationService.ParseToken(cookie.Value)
	if err != nil {
		return nil, err
	}

	return userID, nil
}

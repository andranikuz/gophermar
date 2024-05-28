package handler

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type userBalanceResponse struct {
	current   float64 `json:"current"`
	withdrawn float64 `json:"withdrawn"`
}

func (h HTTPHandler) UserBalance(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := h.GetUserID(r)
	balance, err := h.transactionService.UserBalance(ctx, userID)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(userBalanceResponse{
		current:   balance.Current,
		withdrawn: balance.Withdrawn,
	})
	if err != nil {
		log.Error().Msg(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(resp); err != nil {
		log.Error().Msg(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
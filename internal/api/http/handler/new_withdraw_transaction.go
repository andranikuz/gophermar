package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"github.com/theplant/luhn"
	"gopkg.in/go-playground/validator.v9"
)

type newWithdrawTransactionRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum" validate:"min=1"`
}

func (h HTTPHandler) NewWithdrawTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req newWithdrawTransactionRequest
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
	orderNum, err := strconv.Atoi(req.Order)
	if err != nil {
		log.Info().Msg(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !luhn.Valid(orderNum) {
		log.Info().Msg(`not valid number`)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	userID, _ := h.GetUserID(r)
	id, _ := uuid.NewV6()
	err = h.transactionService.NewTransaction(ctx, id, orderNum, userID, req.Sum)
	if err != nil {
		log.Error().Msg(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

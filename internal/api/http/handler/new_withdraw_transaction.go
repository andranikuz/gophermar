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

	transaction2 "github.com/andranikuz/gophermart/pkg/domain/transaction"
)

type newWithdrawTransactionRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum" validate:"min=1"`
}

func (h HTTPHandler) NewWithdrawTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req newWithdrawTransactionRequest
	body, err := io.ReadAll(r.Body)
	defer logErrorIfExists(err)
	if err = json.Unmarshal(body, &req); err != nil {
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
	orderNum, validationErr := strconv.Atoi(req.Order)
	if validationErr != nil {
		log.Info().Msg(validationErr.Error())
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
	err = h.transactionService.NewTransaction(ctx, id, orderNum, transaction2.TransactionTypeWithdrawal, userID, req.Sum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

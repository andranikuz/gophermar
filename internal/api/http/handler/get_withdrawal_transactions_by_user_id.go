package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/andranikuz/gophermart/pkg/domain/transaction"
)

type withdrawalTransactionResponse []withdrawalTransaction

type withdrawalTransaction struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func (h HTTPHandler) WithdrawalTransactionsByUserID(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := h.GetUserID(r)
	orders, err := h.transactionService.UserTransactionsByType(ctx, userID, transaction.TransactionTypeWithdrawal)
	if err != nil {
		log.Error().Msg(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var response withdrawalTransactionResponse
	for _, o := range orders {
		response = append(
			response,
			withdrawalTransaction{
				Order:       strconv.Itoa(o.OrderNumber),
				Sum:         o.Amount,
				ProcessedAt: o.CreatedAt.Format(time.RFC3339),
			},
		)
	}
	resp, err := json.Marshal(response)
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

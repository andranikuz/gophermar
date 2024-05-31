package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"github.com/theplant/luhn"

	"github.com/andranikuz/gophermart/internal/api"
	"github.com/andranikuz/gophermart/internal/services/order"
)

func (h HTTPHandler) SetOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Info().Msg(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orderNum, err := strconv.Atoi(string(body))
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
	err = h.orderService.SetOrder(ctx, id, orderNum, userID)
	if err != nil {
		if errors.Is(err, order.ErrAccrualTransactionCreatedByAnotherUser) {
			log.Info().Msg(err.Error())
			w.WriteHeader(http.StatusConflict)
			return
		} else if errors.Is(err, order.ErrAccrualTransactionCreatedBySameUser) {
			log.Info().Msg(err.Error())
			w.WriteHeader(http.StatusOK)
			return
		} else {
			log.Error().Msg(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	h.accrualClient.ProcessOrder(api.OrderJob{
		CTX:    ctx,
		Number: orderNum,
		UserID: userID,
	})

	w.WriteHeader(http.StatusAccepted)
}

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
	defer logErrorIfExists(err)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orderNum, err := strconv.Atoi(string(body))
	if err != nil {
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
	setOrderErr := h.orderService.SetOrder(ctx, id, orderNum, userID)
	if setOrderErr != nil {
		if errors.Is(setOrderErr, order.ErrAccrualTransactionCreatedByAnotherUser) {
			log.Info().Msg(setOrderErr.Error())
			w.WriteHeader(http.StatusConflict)
			return
		} else if errors.Is(setOrderErr, order.ErrAccrualTransactionCreatedBySameUser) {
			log.Info().Msg(setOrderErr.Error())
			w.WriteHeader(http.StatusOK)
			return
		} else {
			logErrorIfExists(setOrderErr)
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

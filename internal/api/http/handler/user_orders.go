package handler

import (
	"context"
	"encoding/json"
	"github.com/andranikuz/gophermart/pkg/domain/order"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type userOrdersResponse []orderItem

type orderItem struct {
	Number     int               `json:"number"`
	Status     order.OrderStatus `json:"status"`
	Accrual    float64           `json:"accrual,omitempty"`
	UploadedAt string            `json:"uploadedAt"`
}

func (h HTTPHandler) UserOrders(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := h.GetUserID(r)
	orders, err := h.orderService.UserOrders(ctx, userID)
	if err != nil {
		log.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var response userOrdersResponse
	for _, o := range orders {
		response = append(
			response,
			orderItem{
				Number:     o.Number,
				Status:     o.Status,
				Accrual:    o.Accrual,
				UploadedAt: o.CreatedAt.Format(time.RFC3339),
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

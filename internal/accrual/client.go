package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"

	"github.com/andranikuz/gophermart/internal/config"
	orderService "github.com/andranikuz/gophermart/internal/services/order"
	"github.com/andranikuz/gophermart/internal/services/transaction"
	"github.com/andranikuz/gophermart/pkg/domain/order"
)

type AccrualResponse struct {
	Order   string
	Status  string
	Accrual float64
}

type AccrualClient struct {
	orderService       *orderService.OrderService
	transactionService *transaction.TransactionService
}

func NewAccrualClient(
	orderService *orderService.OrderService,
	transactionService *transaction.TransactionService,
) *AccrualClient {
	return &AccrualClient{
		orderService:       orderService,
		transactionService: transactionService,
	}
}

func (c *AccrualClient) ProcessOrder(ctx context.Context, number int, userID *uuid.UUID) {
	urlString := fmt.Sprintf("%s/api/orders/%d", config.Config.AccrualSystemAddress, number)
	request, err := http.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		return
		log.Error().Msg(err.Error())
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	defer resp.Body.Close()
	log.Info().Msg(`process order`)
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		resp.Body.Close()
		var response AccrualResponse
		if err := json.Unmarshal(body, &response); err != nil {
			log.Error().Msg(err.Error())
		}
		err = c.orderService.UpdateOrderStatus(ctx, number, order.OrderStatus(response.Status))
		if err != nil {
			log.Error().Msg(err.Error())
		}

		if response.Status == "INVALID" || response.Status == "PROCESSED" {
			id, _ := uuid.NewV6()
			err = c.transactionService.NewTransaction(
				ctx,
				id,
				number,
				userID,
				response.Accrual,
			)
		} else {
			time.Sleep(1 * time.Second)
			c.ProcessOrder(ctx, number, userID)
		}
	}
}

package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"

	"github.com/andranikuz/gophermart/internal/api"
	"github.com/andranikuz/gophermart/internal/config"
	"github.com/andranikuz/gophermart/pkg/domain/order"
	transactionDomain "github.com/andranikuz/gophermart/pkg/domain/transaction"
)

type AccrualClient struct {
	orderService       api.OrderServiceInterface
	transactionService api.TransactionServiceInterface
	ch                 chan api.OrderJob
}

var errTooManyRequests = errors.New("accrual: too many requests")
var errNotFinalStatus = errors.New("accrual: not final status")

type response struct {
	Order   string
	Status  string
	Accrual float64
}

func NewAccrualClient(
	orderService api.OrderServiceInterface,
	transactionService api.TransactionServiceInterface,
) *AccrualClient {
	return &AccrualClient{
		orderService:       orderService,
		transactionService: transactionService,
		ch:                 make(chan api.OrderJob, 100),
	}
}

// Worker запускаем воркер
func (c *AccrualClient) Worker() {
	for {
		job := <-c.ch
		err := c.do(job)
		if err != nil {
			if errors.Is(err, errTooManyRequests) {
				log.Info().Msg(err.Error())
				time.Sleep(60 * time.Second)
				c.ProcessOrder(job)
			} else if errors.Is(err, errNotFinalStatus) {
				log.Info().Msg(err.Error())
				c.ProcessOrder(job)
			} else {
				log.Error().Msg(err.Error())
				time.Sleep(10 * time.Second)
				c.ProcessOrder(job)
			}
		}
	}
}

// ProcessOrder помещаем заказ в очередь
func (c AccrualClient) ProcessOrder(job api.OrderJob) {
	c.ch <- job
}

// обрабатываем заказ из очереди
func (c *AccrualClient) do(job api.OrderJob) error {
	urlString := fmt.Sprintf("%s/api/orders/%d", config.Config.AccrualSystemAddress, job.Number)
	request, err := http.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Info().Msg(`accrual: process order`)
	if resp.StatusCode == http.StatusTooManyRequests {
		return errTooManyRequests
	}
	if resp.StatusCode != http.StatusOK {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var r response
	if err = json.Unmarshal(body, &r); err != nil {
		return err
	}
	err = c.orderService.UpdateOrderStatus(job.CTX, job.Number, order.OrderStatus(r.Status))
	if r.Status != "INVALID" && r.Status != "PROCESSED" {
		return errNotFinalStatus
	}
	id, _ := uuid.NewV6()
	return c.transactionService.NewTransaction(
		job.CTX,
		id,
		job.Number,
		transactionDomain.TransactionTypeAccrual,
		job.UserID,
		r.Accrual,
	)
}

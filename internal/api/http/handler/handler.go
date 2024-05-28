package handler

import (
	"context"
	"github.com/andranikuz/gophermart/internal/accrual"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/andranikuz/gophermart/internal/container"
	"github.com/andranikuz/gophermart/internal/services/auth"
	"github.com/andranikuz/gophermart/internal/services/order"
	"github.com/andranikuz/gophermart/internal/services/transaction"
)

type HTTPHandler struct {
	authenticationService *auth.AuthenticationService
	transactionService    *transaction.TransactionService
	orderService          *order.OrderService
	accrualClient         *accrual.AccrualClient
}

func NewHTTPHandler(cnt *container.Container) HTTPHandler {
	h := HTTPHandler{
		authenticationService: cnt.AuthenticationService(),
		transactionService:    cnt.TransactionService(),
		orderService:          cnt.OrderService(),
		accrualClient:         cnt.AccrualClient(),
	}

	return h
}

func (h HTTPHandler) Router(ctx context.Context) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer)
	r.Post("/api/user/register", func(w http.ResponseWriter, r *http.Request) {
		h.RegisterHandler(ctx, w, r)
	})
	r.Post("/api/user/login", func(w http.ResponseWriter, r *http.Request) {
		h.LoginHandler(ctx, w, r)
	})

	r.Group(func(r chi.Router) {
		r.Use(h.authMiddleware)
		r.Post("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
			h.SetOrder(ctx, w, r)
		})
		r.Get("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
			h.UserOrders(ctx, w, r)
		})
		r.Get("/api/user/balance", func(w http.ResponseWriter, r *http.Request) {
			h.UserBalance(ctx, w, r)
		})
		r.Post("/api/user/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {
			h.NewWithdrawTransaction(ctx, w, r)
		})
		r.Get("/api/user/withdrawals", func(w http.ResponseWriter, r *http.Request) {
			h.WithdrawalTransactionsByUserID(ctx, w, r)
		})
	})

	return r
}

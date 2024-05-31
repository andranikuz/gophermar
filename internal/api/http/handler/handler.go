package handler

import (
	"context"
	"github.com/andranikuz/gophermart/internal/api"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/andranikuz/gophermart/internal/container"
)

type HTTPHandler struct {
	authenticationService api.AuthenticationServiceInterface
	transactionService    api.TransactionServiceInterface
	orderService          api.OrderServiceInterface
	accrualClient         api.AccrualClientInterface
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

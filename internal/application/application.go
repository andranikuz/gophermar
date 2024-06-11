package application

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andranikuz/gophermart/internal/api/http/handler"
	"github.com/andranikuz/gophermart/internal/config"
	"github.com/andranikuz/gophermart/internal/container"
	"github.com/andranikuz/gophermart/internal/postgres"
)

type Application struct {
	cnt *container.Container
	ctx context.Context
}

func NewApplication() (*Application, error) {
	config.Init()
	pool, err := pgxpool.New(context.Background(), config.Config.DatabaseDSN)
	ctx := context.Background()
	if err != nil {
		return nil, err
	}
	userRepo := postgres.NewUserRepository(pool)
	err = userRepo.CreateTable(ctx)
	if err != nil {
		return nil, err
	}
	transactionRepo := postgres.NewTransactionRepository(pool)
	err = transactionRepo.CreateTable(ctx)
	if err != nil {
		return nil, err
	}
	orderRepo := postgres.NewOrderRepositoryRepository(pool)
	err = orderRepo.CreateTable(ctx)
	if err != nil {
		return nil, err
	}
	cnt := container.NewContainer(userRepo, transactionRepo, orderRepo)

	return &Application{
		ctx: ctx,
		cnt: cnt,
	}, nil
}

func (a *Application) Run() error {
	a.runWorker()
	return a.runServer()
}

func (a Application) runServer() error {
	httpHandler := handler.NewHTTPHandler(a.cnt)
	return http.ListenAndServe(config.Config.ServerAddress, httpHandler.Router(a.ctx))
}

func (a Application) runWorker() {
	go a.cnt.AccrualClient().Worker()
}

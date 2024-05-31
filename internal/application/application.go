package application

import (
	"context"
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

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
	db, err := sql.Open("pgx", config.Config.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	userRepo := postgres.NewUserRepository(db)
	err = userRepo.CreateTable()
	if err != nil {
		return nil, err
	}
	transactionRepo := postgres.NewTransactionRepository(db)
	err = transactionRepo.CreateTable()
	if err != nil {
		return nil, err
	}
	orderRepo := postgres.NewOrderRepositoryRepository(db)
	err = orderRepo.CreateTable()
	if err != nil {
		return nil, err
	}
	cnt := container.NewContainer(userRepo, transactionRepo, orderRepo)

	return &Application{
		ctx: context.Background(),
		cnt: cnt,
	}, nil
}

func (a *Application) Run() error {
	httpHandler := handler.NewHTTPHandler(a.cnt)
	return http.ListenAndServe(config.Config.ServerAddress, httpHandler.Router(a.ctx))
}

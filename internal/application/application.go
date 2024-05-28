package application

import (
	"context"
	"net/http"

	"github.com/andranikuz/gophermart/internal/api/http/handler"
	"github.com/andranikuz/gophermart/internal/config"
	"github.com/andranikuz/gophermart/internal/container"
)

type Application struct {
	cnt *container.Container
	ctx context.Context
}

func NewApplication() (*Application, error) {
	config.Init()
	cnt := container.NewContainer()
	a := Application{
		ctx: context.Background(),
		cnt: cnt,
	}

	return &a, nil
}

func (a *Application) Run() error {
	httpHandler := handler.NewHTTPHandler(a.cnt)
	return http.ListenAndServe(config.Config.ServerAddress, httpHandler.Router(a.ctx))
}

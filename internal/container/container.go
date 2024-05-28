package container

import (
	"database/sql"
	"github.com/andranikuz/gophermart/internal/accrual"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"

	"github.com/andranikuz/gophermart/internal/config"
	"github.com/andranikuz/gophermart/internal/services/auth"
	service "github.com/andranikuz/gophermart/internal/services/order"
	"github.com/andranikuz/gophermart/internal/services/transaction"
	"github.com/andranikuz/gophermart/pkg/domain/order"
	domain "github.com/andranikuz/gophermart/pkg/domain/transaction"
	"github.com/andranikuz/gophermart/pkg/domain/user"
)

type Container struct {
	db *sql.DB
	// repositories
	userRepo        user.Repository
	transactionRepo domain.Repository
	orderRepo       order.Repository
	//services
	authenticationService *auth.AuthenticationService
	transactionService    *transaction.TransactionService
	orderService          *service.OrderService
	// client
	accrualClient *accrual.AccrualClient
}

func NewContainer() *Container {
	c := Container{}
	c.UserRepository()
	c.TransactionRepository()
	c.OrderRepository()

	return &c
}

func (c *Container) DB() *sql.DB {
	if c.db == nil {
		db, err := sql.Open("pgx", config.Config.DatabaseDSN)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		return db
	}
	return c.db
}

func (c *Container) AuthenticationService() *auth.AuthenticationService {
	if c.authenticationService == nil {
		c.authenticationService = auth.NewAuthService(c.UserRepository())
	}
	return c.authenticationService
}

func (c *Container) TransactionService() *transaction.TransactionService {
	if c.transactionService == nil {
		c.transactionService = transaction.NewTransactionService(c.TransactionRepository())
	}
	return c.transactionService
}

func (c *Container) OrderService() *service.OrderService {
	if c.orderService == nil {
		c.orderService = service.NewOrderService(c.OrderRepository())
	}
	return c.orderService
}

func (c *Container) AccrualClient() *accrual.AccrualClient {
	if c.accrualClient == nil {
		c.accrualClient = accrual.NewAccrualClient(
			c.OrderService(),
			c.TransactionService(),
		)
	}
	return c.accrualClient
}

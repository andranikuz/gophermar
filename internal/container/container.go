package container

import (
	"github.com/andranikuz/gophermart/internal/accrual"
	"github.com/andranikuz/gophermart/internal/api"
	"github.com/andranikuz/gophermart/internal/services/auth"
	service "github.com/andranikuz/gophermart/internal/services/order"
	"github.com/andranikuz/gophermart/internal/services/transaction"
	"github.com/andranikuz/gophermart/pkg/domain/order"
	domain "github.com/andranikuz/gophermart/pkg/domain/transaction"
	"github.com/andranikuz/gophermart/pkg/domain/user"
)

type Container struct {
	// repositories
	userRepo        user.Repository
	transactionRepo domain.Repository
	orderRepo       order.Repository
	//services
	authenticationService api.AuthenticationServiceInterface
	transactionService    api.TransactionServiceInterface
	orderService          api.OrderServiceInterface
	// client
	accrualClient api.AccrualClientInterface
}

func NewContainer(
	userRepo user.Repository,
	transactionRepo domain.Repository,
	orderRepo order.Repository,
) *Container {
	return &Container{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		orderRepo:       orderRepo,
	}
}

func (c *Container) AuthenticationService() api.AuthenticationServiceInterface {
	if c.authenticationService == nil {
		c.authenticationService = auth.NewAuthService(c.UserRepository())
	}
	return c.authenticationService
}

func (c *Container) TransactionService() api.TransactionServiceInterface {
	if c.transactionService == nil {
		c.transactionService = transaction.NewTransactionService(c.TransactionRepository())
	}
	return c.transactionService
}

func (c *Container) OrderService() api.OrderServiceInterface {
	if c.orderService == nil {
		c.orderService = service.NewOrderService(c.OrderRepository())
	}
	return c.orderService
}

func (c *Container) AccrualClient() api.AccrualClientInterface {
	if c.accrualClient == nil {
		c.accrualClient = accrual.NewAccrualClient(
			c.OrderService(),
			c.TransactionService(),
		)
	}
	return c.accrualClient
}

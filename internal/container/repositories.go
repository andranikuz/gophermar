package container

import (
	"github.com/andranikuz/gophermart/internal/postgres"
	"github.com/andranikuz/gophermart/pkg/domain/order"
	"github.com/andranikuz/gophermart/pkg/domain/transaction"
	"github.com/andranikuz/gophermart/pkg/domain/user"
)

func (c *Container) UserRepository() user.Repository {
	if c.userRepo == nil {
		repo := postgres.NewUserRepository(c.DB())
		repo.CreateTable()
		c.userRepo = repo
	}
	return c.userRepo
}

func (c *Container) TransactionRepository() transaction.Repository {
	if c.transactionRepo == nil {
		repo := postgres.NewTransactionRepository(c.DB())
		repo.CreateTable()
		c.transactionRepo = repo
	}
	return c.transactionRepo
}

func (c *Container) OrderRepository() order.Repository {
	if c.orderRepo == nil {
		repo := postgres.NewOrderRepositoryRepository(c.DB())
		repo.CreateTable()
		c.orderRepo = repo
	}
	return c.orderRepo
}

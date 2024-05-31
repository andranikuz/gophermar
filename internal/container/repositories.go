package container

import (
	"github.com/andranikuz/gophermart/pkg/domain/order"
	"github.com/andranikuz/gophermart/pkg/domain/transaction"
	"github.com/andranikuz/gophermart/pkg/domain/user"
)

func (c *Container) UserRepository() user.Repository {
	return c.userRepo
}

func (c *Container) TransactionRepository() transaction.Repository {
	return c.transactionRepo
}

func (c *Container) OrderRepository() order.Repository {
	return c.orderRepo
}

package order

import (
	"time"

	"github.com/gofrs/uuid"
)

type Order struct {
	ID        uuid.UUID
	Number    int
	UserID    *uuid.UUID
	Status    OrderStatus
	Accrual   float64
	CreatedAt time.Time
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

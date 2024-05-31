package api

import (
	"context"

	"github.com/gofrs/uuid"
)

type AccrualClientInterface interface {
	ProcessOrder(ctx context.Context, number int, userID *uuid.UUID)
}

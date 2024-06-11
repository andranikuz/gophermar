package api

import (
	"context"

	"github.com/gofrs/uuid"
)

type OrderJob struct {
	CTX    context.Context
	Number int
	UserID *uuid.UUID
}

type AccrualClientInterface interface {
	ProcessOrder(job OrderJob)
	Worker()
}

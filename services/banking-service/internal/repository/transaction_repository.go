package repository

import (
	"banking-service/internal/model"
	"context"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *model.Transaction) error
}

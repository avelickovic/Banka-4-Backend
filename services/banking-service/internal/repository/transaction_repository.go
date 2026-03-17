package repository

import (
	"banking-service/internal/model"
	"context"
	
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *model.Transaction) error
	GetByID(ctx context.Context, transactionID uint) (*model.Transaction, error)

	GetByIDWithTx(ctx context.Context, db *gorm.DB, transactionID uint) (*model.Transaction, error)
	UpdateWithTx(ctx context.Context, db *gorm.DB, transaction *model.Transaction) error
}

package repository

import (
	"banking-service/internal/model"
	"context"
	
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	AccountNumberExists(ctx context.Context, accountNumber string) (bool, error)

	GetByAccountNumberWithTx(ctx context.Context, db *gorm.DB, accountNumber string) (*model.Account, error)
	UpdateWithTx(ctx context.Context, db *gorm.DB, account *model.Account) error
}

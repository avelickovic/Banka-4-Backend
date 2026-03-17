package repository

import (
	"banking-service/internal/model"
	"context"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *model.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) GetByID(ctx context.Context, id uint) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.WithContext(ctx).First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetByIDWithTx(ctx context.Context, db *gorm.DB, id uint) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := db.WithContext(ctx).First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}


func (r *transactionRepository) UpdateWithTx(ctx context.Context, db *gorm.DB, transaction *model.Transaction) error {
	return db.WithContext(ctx).Save(transaction).Error
}

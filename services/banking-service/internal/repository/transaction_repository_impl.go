package repository

import (
	"banking-service/internal/model"
	"common/pkg/db"
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
	currentDB := db.DBFromContext(ctx, r.db)
	return currentDB.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) GetByID(ctx context.Context, id uint) (*model.Transaction, error) {
	db := db.DBFromContext(ctx, r.db)

	var transaction model.Transaction
	if err := db.WithContext(ctx).First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetByPayerAccountNumber(ctx context.Context, accountNumber string) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	err := r.db.WithContext(ctx).Where("payer_account_number = ?", accountNumber).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetByRecipientAccountNumber(ctx context.Context, accountNumber string) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	err := r.db.WithContext(ctx).Where("recipient_account_number = ?", accountNumber).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) Update(ctx context.Context, transaction *model.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

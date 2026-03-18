package repository

import (
	"banking-service/internal/model"
	"context"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) GetByID(ctx context.Context, id uint) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.WithContext(ctx).Preload("Transaction").First(&payment, id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *paymentRepository) FindAllByClientID(ctx context.Context, clientID uint, filter PaymentFilter) ([]model.Payment, error) {
	var payments []model.Payment

	query := r.db.WithContext(ctx).
		Joins("JOIN transactions ON transactions.transaction_id = payments.transaction_id").
		Joins("JOIN accounts ON accounts.account_number = transactions.payer_account_number").
		Where("accounts.client_id = ?", clientID).
		Preload("Transaction")

	if filter.DateFrom != nil {
		query = query.Where("transactions.created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("transactions.created_at <= ?", filter.DateTo)
	}
	if filter.AmountMin != nil {
		query = query.Where("transactions.start_amount >= ?", filter.AmountMin)
	}
	if filter.AmountMax != nil {
		query = query.Where("transactions.start_amount <= ?", filter.AmountMax)
	}
	if filter.Status != nil {
		query = query.Where("transactions.status = ?", filter.Status)
	}

	err := query.Find(&payments).Error
	return payments, err
}

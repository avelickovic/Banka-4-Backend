package repository

import (
	"banking-service/internal/model"
	"common/pkg/db"
	"context"

	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *model.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *accountRepository) AccountNumberExists(ctx context.Context, accountNumber string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.Account{}).
		Where("account_number = ?", accountNumber).
		Count(&count).
		Error

	return count > 0, err
}

func (r *accountRepository) FindAllByClientID(ctx context.Context, clientID uint) ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.WithContext(ctx).
		Preload("Currency").
		Where("client_id = ? AND status = ?", clientID, "Active").
		Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) FindByAccountNumber(ctx context.Context, accountNumber string) (*model.Account, error) {
	var account model.Account
	err := r.db.WithContext(ctx).
		Preload("Currency").
		Where("account_number = ?", accountNumber).
		First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindByAccountNumberAndClientID(ctx context.Context, accountNumber string, clientID uint) (*model.Account, error) {
	var account model.Account
	err := r.db.WithContext(ctx).
		Preload("Currency").
		Where("account_number = ? AND client_id = ?", accountNumber, clientID).
		First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) UpdateName(ctx context.Context, accountNumber string, name string) error {
	return r.db.WithContext(ctx).
		Model(&model.Account{}).
		Where("account_number = ?", accountNumber).
		Update("name", name).Error
}

func (r *accountRepository) UpdateLimits(ctx context.Context, accountNumber string, daily float64, monthly float64) error {
	return r.db.WithContext(ctx).
		Model(&model.Account{}).
		Where("account_number = ?", accountNumber).
		Updates(map[string]interface{}{"daily_limit": daily, "monthly_limit": monthly}).Error
}

func (r *accountRepository) NameExistsForClient(ctx context.Context, clientID uint, name string, excludeNumber string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Account{}).
		Where("client_id = ? AND name = ? AND account_number != ?", clientID, name, excludeNumber).
		Count(&count).Error
	return count > 0, err
}

func (r *accountRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*model.Account, error) {
	db := db.DBFromContext(ctx, r.db)

	var account model.Account
	if err := db.WithContext(ctx).First(&account, accountNumber).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) Update(ctx context.Context, account *model.Account) error {
	db := db.DBFromContext(ctx, r.db)

	return db.WithContext(ctx).Save(account).Error
}


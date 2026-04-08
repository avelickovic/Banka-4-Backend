package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type otcContractRepositoryImpl struct {
	db *gorm.DB
}

func NewOtcContractRepository(db *gorm.DB) OtcContractRepository {
	return &otcContractRepositoryImpl{db: db}
}

func (r *otcContractRepositoryImpl) Create(ctx context.Context, contract *model.OtcContract) error {
	return r.db.WithContext(ctx).Create(contract).Error
}

func (r *otcContractRepositoryImpl) FindByID(ctx context.Context, id uint) (*model.OtcContract, error) {
	var contract model.OtcContract
	result := r.db.WithContext(ctx).Preload("Asset").First(&contract, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &contract, result.Error
}

func (r *otcContractRepositoryImpl) Save(ctx context.Context, contract *model.OtcContract) error {
	return r.db.WithContext(ctx).Save(contract).Error
}

func (r *otcContractRepositoryImpl) FindByBuyerID(ctx context.Context, buyerID uint) ([]model.OtcContract, error) {
	var contracts []model.OtcContract
	err := r.db.WithContext(ctx).Preload("Asset").Where("buyer_id = ?", buyerID).Find(&contracts).Error
	return contracts, err
}

func (r *otcContractRepositoryImpl) FindBySellerID(ctx context.Context, sellerID uint) ([]model.OtcContract, error) {
	var contracts []model.OtcContract
	err := r.db.WithContext(ctx).Preload("Asset").Where("seller_id = ?", sellerID).Find(&contracts).Error
	return contracts, err
}

func (r *otcContractRepositoryImpl) FindPendingBankApproval(ctx context.Context) ([]model.OtcContract, error) {
	var contracts []model.OtcContract
	err := r.db.WithContext(ctx).Preload("Asset").Where("bank_approved IS NULL").Find(&contracts).Error
	return contracts, err
}
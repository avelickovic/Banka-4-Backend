package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type OtcContractRepository interface {
	Create(ctx context.Context, contract *model.OtcContract) error
	FindByID(ctx context.Context, id uint) (*model.OtcContract, error)
	Save(ctx context.Context, contract *model.OtcContract) error
	FindByBuyerID(ctx context.Context, buyerID uint) ([]model.OtcContract, error)
	FindBySellerID(ctx context.Context, sellerID uint) ([]model.OtcContract, error)
	FindPendingBankApproval(ctx context.Context) ([]model.OtcContract, error)
}
package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type OtcShareReservationRepository interface {
	Create(ctx context.Context, reservation *model.OtcShareReservation) error
	FindByContractID(ctx context.Context, contractID uint) (*model.OtcShareReservation, error)
	FindByContractIDForUpdate(ctx context.Context, contractID uint) (*model.OtcShareReservation, error)

	// SumActiveReservedBySellerAsset sums ReservedAmount across ACTIVE OTC share
	// reservations for the given seller/owner type/stock, excluding the provided
	// contract when excludeContractID is not nil.
	SumActiveReservedBySellerAsset(ctx context.Context, sellerID uint, ownerType model.OwnerType, stockAssetID uint, excludeContractID *uint) (float64, error)
	Save(ctx context.Context, reservation *model.OtcShareReservation) error
}

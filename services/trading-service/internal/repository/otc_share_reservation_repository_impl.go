package repository

import (
	"context"
	"errors"

	commondb "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/db"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type otcShareReservationRepository struct {
	db *gorm.DB
}

func NewOtcShareReservationRepository(db *gorm.DB) OtcShareReservationRepository {
	return &otcShareReservationRepository{db: db}
}

func (r *otcShareReservationRepository) Create(ctx context.Context, reservation *model.OtcShareReservation) error {
	return commondb.DBFromContext(ctx, r.db).Create(reservation).Error
}

func (r *otcShareReservationRepository) FindByContractID(ctx context.Context, contractID uint) (*model.OtcShareReservation, error) {
	return r.findByContractID(ctx, contractID, false)
}

func (r *otcShareReservationRepository) FindByContractIDForUpdate(ctx context.Context, contractID uint) (*model.OtcShareReservation, error) {
	return r.findByContractID(ctx, contractID, true)
}

func (r *otcShareReservationRepository) SumActiveReservedBySellerAsset(ctx context.Context, sellerID uint, ownerType model.OwnerType, stockAssetID uint, excludeContractID *uint) (float64, error) {
	query := commondb.DBFromContext(ctx, r.db).
		Model(&model.OtcShareReservation{}).
		Where(
			"seller_id = ? AND owner_type = ? AND stock_asset_id = ? AND status = ?",
			sellerID,
			ownerType,
			stockAssetID,
			model.OtcShareReservationStatusActive,
		)

	if excludeContractID != nil {
		query = query.Where("contract_id <> ?", *excludeContractID)
	}

	var total float64
	if err := query.Select("COALESCE(SUM(reserved_amount), 0)").Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *otcShareReservationRepository) Save(ctx context.Context, reservation *model.OtcShareReservation) error {
	return commondb.DBFromContext(ctx, r.db).Save(reservation).Error
}

func (r *otcShareReservationRepository) findByContractID(ctx context.Context, contractID uint, forUpdate bool) (*model.OtcShareReservation, error) {
	query := commondb.DBFromContext(ctx, r.db).Where("contract_id = ?", contractID)
	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var reservation model.OtcShareReservation
	if err := query.First(&reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &reservation, nil
}

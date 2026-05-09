package repository

import (
	"context"
	stderrors "errors"

	commondb "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type otcOfferRepositoryImpl struct {
	db *gorm.DB
}

func NewOtcOfferRepository(db *gorm.DB) OtcOfferRepository {
	return &otcOfferRepositoryImpl{db: db}
}

func (r *otcOfferRepositoryImpl) Create(ctx context.Context, offer *model.OtcOffer) error {
	return commondb.DBFromContext(ctx, r.db).Create(offer).Error
}

func (r *otcOfferRepositoryImpl) Save(ctx context.Context, offer *model.OtcOffer) error {
	return commondb.DBFromContext(ctx, r.db).Save(offer).Error
}

func (r *otcOfferRepositoryImpl) FindByID(ctx context.Context, id uint) (*model.OtcOffer, error) {
	return r.findByID(ctx, id, false)
}

func (r *otcOfferRepositoryImpl) FindByIDForUpdate(ctx context.Context, id uint) (*model.OtcOffer, error) {
	return r.findByID(ctx, id, true)
}

func (r *otcOfferRepositoryImpl) FindActiveForUser(ctx context.Context, userID uint) ([]model.OtcOffer, error) {
	var offers []model.OtcOffer
	err := commondb.DBFromContext(ctx, r.db).
		Preload("Stock").
		Preload("Stock.Asset").
		Where("(buyer_id = ? OR seller_id = ?) AND status = ?", userID, userID, model.OtcOfferStatusActive).
		Order("last_modified DESC").
		Find(&offers).Error
	return offers, err
}

func (r *otcOfferRepositoryImpl) FindActiveBySellerAndStock(ctx context.Context, sellerID, stockAssetID uint, excludeOfferID *uint) ([]model.OtcOffer, error) {
	var offers []model.OtcOffer
	q := commondb.DBFromContext(ctx, r.db).
		Where("seller_id = ? AND stock_asset_id = ? AND status = ?", sellerID, stockAssetID, model.OtcOfferStatusActive)
	if excludeOfferID != nil {
		q = q.Where("otc_offer_id <> ?", *excludeOfferID)
	}
	err := q.Find(&offers).Error
	return offers, err
}

func (r *otcOfferRepositoryImpl) findByID(ctx context.Context, id uint, forUpdate bool) (*model.OtcOffer, error) {
	query := commondb.DBFromContext(ctx, r.db).
		Preload("Stock").
		Preload("Stock.Asset")

	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var offer model.OtcOffer
	result := query.First(&offer, id)

	if stderrors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}
	return &offer, nil
}

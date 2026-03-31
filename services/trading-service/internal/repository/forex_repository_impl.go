package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

type forexRepository struct {
	db *gorm.DB
}

func NewForexRepository(db *gorm.DB) ForexRepository {
	return &forexRepository{db: db}
}

func (r *forexRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.ForexPair{}).
		Count(&count).Error

	return count, err
}

func (r *forexRepository) Upsert(ctx context.Context, pair model.ForexPair) error {
	return r.db.WithContext(ctx).
		Where("base = ? AND quote = ?", pair.Base, pair.Quote).
		Assign(pair).
		FirstOrCreate(&pair).Error
}

func (r *forexRepository) FindByListingIDs(listingIDs []uint) ([]model.ForexPair, error) {
	var pairs []model.ForexPair
	if err := r.db.Where("listing_id IN ?", listingIDs).Preload("Listing").Find(&pairs).Error; err != nil {
		return nil, err
	}
	return pairs, nil
}

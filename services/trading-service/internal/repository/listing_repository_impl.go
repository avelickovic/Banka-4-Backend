package repository

import (
	"context"
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

type listingRepository struct {
	db *gorm.DB
}

func NewListingRepository(db *gorm.DB) ListingRepository {
	return &listingRepository{db: db}
}

func (r *listingRepository) FindAll() ([]model.Listing, error) {
	var listings []model.Listing
	if err := r.db.Find(&listings).Error; err != nil {
		return nil, err
	}
	return listings, nil
}

func (r *listingRepository) Upsert(listing *model.Listing) error {
	return r.db.
		Where(model.Listing{Ticker: listing.Ticker}).
		Assign(*listing).
		FirstOrCreate(listing).Error
}

func (r *listingRepository) UpdatePriceAndAsk(listing *model.Listing, price, ask float64) error {
	return r.db.Model(listing).Updates(map[string]interface{}{
		"price":        price,
		"ask":          ask,
		"last_refresh": time.Now(),
	}).Error
}

func (r *listingRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Listing{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

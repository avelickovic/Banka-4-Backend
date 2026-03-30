package repository

import (
	"context"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db: db}
}

// Upsert inserts or updates a Stock matched by ListingID.
func (r *stockRepository) Upsert(ctx context.Context, stock *model.Stock) error {
	return r.db.WithContext(ctx).
		Where(model.Stock{ListingID: stock.ListingID}).
		Assign(*stock).
		FirstOrCreate(stock).Error
}

func (r *stockRepository) FindAll(ctx context.Context) ([]model.Stock, error) {
	var stocks []model.Stock
	if err := r.db.WithContext(ctx).Preload("Listing").Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}


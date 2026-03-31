package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

type futuresRepository struct {
	db *gorm.DB
}

func NewFuturesRepository(db *gorm.DB) FuturesRepository {
	return &futuresRepository{db: db}
}

func (r *futuresRepository) FindByTickers(ctx context.Context, tickers []string) ([]model.FuturesContract, error) {
	var contracts []model.FuturesContract
	err := r.db.WithContext(ctx).
		Where("ticker IN ?", tickers).
		Find(&contracts).Error
	return contracts, err
}

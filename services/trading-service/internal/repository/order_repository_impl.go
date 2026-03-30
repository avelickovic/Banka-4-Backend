package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepositoryImpl{db: db}
}

func (r *orderRepositoryImpl) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepositoryImpl) FindByID(ctx context.Context, id uint) (*model.Order, error) {
	var order model.Order
	result := r.db.WithContext(ctx).Preload("Listing").First(&order, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &order, result.Error
}

func (r *orderRepositoryImpl) Save(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

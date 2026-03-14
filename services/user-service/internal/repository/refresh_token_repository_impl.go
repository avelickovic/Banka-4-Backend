package repository

import (
	"context"
	"errors"
	"user-service/internal/model"

	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var t model.RefreshToken
	result := r.db.WithContext(ctx).Where("token = ?", token).First(&t)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &t, nil
}

func (r *refreshTokenRepository) DeleteByIdentityID(ctx context.Context, identityID uint) error {
	return r.db.WithContext(ctx).Where("identity_id = ?", identityID).Delete(&model.RefreshToken{}).Error
}

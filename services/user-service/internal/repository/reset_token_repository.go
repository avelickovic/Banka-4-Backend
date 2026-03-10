package repository

import (
	"context"
	"user-service/internal/model"
)

type ResetTokenRepository interface {
	Create(ctx context.Context, token *model.ResetToken) error
	FindByEmployeeID(ctx context.Context, employeeID uint) (*model.ResetToken, error)
	FindByCode(ctx context.Context, code string) (*model.ResetToken, error)
	DeleteByEmployeeID(ctx context.Context, employeeID uint) error
}

package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type FuturesRepository interface {
	FindByTickers(ctx context.Context, tickers []string) ([]model.FuturesContract, error)
}

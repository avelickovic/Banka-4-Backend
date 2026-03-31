package repository

import "github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"

type OrderOwnershipRepository interface {
	FindByIdentity(identityID uint, ownerType model.OwnerType) ([]model.OrderOwnership, error)
}

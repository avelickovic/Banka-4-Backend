package dto

import (
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/permission"
)

type UpdateClientRequest struct {
	FirstName   *string                  `json:"first_name" binding:"omitempty,max=20"`
	LastName    *string                  `json:"last_name" binding:"omitempty,max=100"`
	Gender      *string                  `json:"gender"`
	DateOfBirth *time.Time               `json:"date_of_birth"`
	Email       *string                  `json:"email" binding:"omitempty,email"`
	PhoneNumber *string                  `json:"phone_number"`
	Address     *string                  `json:"address"`
	Username    *string                  `json:"username"`
	Permissions *[]permission.Permission `json:"permissions" binding:"omitempty,unique_permissions,dive,permission"`
}

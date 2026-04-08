package model

import "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/permission"

type EmployeePermission struct {
	EmployeeID uint                  `gorm:"primaryKey"`
	Permission permission.Permission `gorm:"type:varchar(64);primaryKey"`
}

func (EmployeePermission) TableName() string {
	return "employee_permissions"
}

type ClientPermission struct {
	ClientID   uint                  `gorm:"primaryKey"`
	Permission permission.Permission `gorm:"type:varchar(64);primaryKey"`
}

func (ClientPermission) TableName() string {
	return "client_permissions"
}

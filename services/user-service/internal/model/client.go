package model

import (
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/permission"
)

type Client struct {
	ClientID                 uint   `gorm:"primaryKey"`
	IdentityID               uint   `gorm:"uniqueIndex;not null"`
	FirstName                string `gorm:"size:20;not null"`
	LastName                 string `gorm:"size:100;not null"`
	MobileVerificationSecret string `gorm:"size:64"`
	DateOfBirth              time.Time
	Gender                   string `gorm:"size:10"`
	PhoneNumber              string `gorm:"size:20"`
	Address                  string `gorm:"size:255"`

	Identity    Identity
	Permissions []ClientPermission `gorm:"foreignKey:ClientID"`
}

func (c *Client) HasPermission(p permission.Permission) bool {
	for _, cp := range c.Permissions {
		if cp.Permission == p {
			return true
		}
	}

	return false
}

func (c *Client) RawPermissions() []permission.Permission {
	if c == nil || len(c.Permissions) == 0 {
		return []permission.Permission{}
	}

	permissions := make([]permission.Permission, 0, len(c.Permissions))
	for _, cp := range c.Permissions {
		permissions = append(permissions, cp.Permission)
	}

	return permissions
}

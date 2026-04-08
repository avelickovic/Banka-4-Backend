package auth

import "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/permission"

func HasPermission(permissions []permission.Permission, required permission.Permission) bool {
	for _, current := range permissions {
		if current == required {
			return true
		}
	}

	return false
}

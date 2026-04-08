package permission

import (
	"context"

	"gorm.io/gorm"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/auth"
	commonjwt "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/jwt"
	perm "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/permission"
)

// DBPermissionProvider loads all permissions for an identity by querying the
// database directly.
type DBPermissionProvider struct {
	db *gorm.DB
}

type permissionLookupConfig struct {
	subjectColumn    string
	permissionsTable string
	permissionColumn string
}

var permissionsByType = map[auth.IdentityType]permissionLookupConfig{
	auth.IdentityEmployee: {
		subjectColumn:    "employee_id",
		permissionsTable: "employee_permissions",
		permissionColumn: "employee_id",
	},
	auth.IdentityClient: {
		subjectColumn:    "client_id",
		permissionsTable: "client_permissions",
		permissionColumn: "client_id",
	},
}

func NewDBPermissionProvider(db *gorm.DB) *DBPermissionProvider {
	return &DBPermissionProvider{db: db}
}

func (p *DBPermissionProvider) GetPermissions(ctx context.Context, claims *commonjwt.Claims) ([]perm.Permission, error) {
	identityType := auth.IdentityType(claims.IdentityType)
	config, ok := permissionsByType[identityType]
	if !ok {
		return []perm.Permission{}, nil
	}

	subjectID, hasSubjectID, err := auth.SubjectIDFromClaims(claims)
	if err != nil {
		return nil, err
	}

	if !hasSubjectID {
		result := p.db.WithContext(ctx).
			Table(identityTable(identityType)).
			Select(config.subjectColumn).
			Where("identity_id = ?", claims.IdentityID).
			Scan(&subjectID)

		if result.Error != nil {
			return nil, result.Error
		}

		if result.RowsAffected == 0 {
			return []perm.Permission{}, nil
		}
	}

	var permissions []perm.Permission
	err = p.db.WithContext(ctx).
		Table(config.permissionsTable).
		Where(config.permissionColumn+" = ?", subjectID).
		Pluck("permission", &permissions).Error

	return permissions, err
}

func identityTable(identityType auth.IdentityType) string {
	switch identityType {
	case auth.IdentityEmployee:
		return "employees"
	case auth.IdentityClient:
		return "clients"
	default:
		return ""
	}
}

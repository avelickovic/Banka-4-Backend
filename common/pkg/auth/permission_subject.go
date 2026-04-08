package auth

import (
	"fmt"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/jwt"
)

func SubjectIDFromClaims(claims *jwt.Claims) (uint, bool, error) {
	if claims == nil {
		return 0, false, fmt.Errorf("claims are required")
	}

	switch IdentityType(claims.IdentityType) {
	case IdentityEmployee:
		if claims.EmployeeID == nil {
			return 0, false, nil
		}

		return *claims.EmployeeID, true, nil
	case IdentityClient:
		if claims.ClientID == nil {
			return 0, false, nil
		}

		return *claims.ClientID, true, nil
	default:
		return 0, false, fmt.Errorf("unsupported identity type: %s", claims.IdentityType)
	}
}

func SetSubjectIDOnClaims(claims *jwt.Claims, identityType IdentityType, subjectID uint) error {
	if claims == nil {
		return fmt.Errorf("claims are required")
	}

	switch identityType {
	case IdentityEmployee:
		claims.EmployeeID = &subjectID
		claims.ClientID = nil
	case IdentityClient:
		claims.ClientID = &subjectID
		claims.EmployeeID = nil
	default:
		return fmt.Errorf("unsupported identity type: %s", identityType)
	}

	return nil
}

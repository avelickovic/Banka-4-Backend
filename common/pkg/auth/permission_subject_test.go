package auth

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/jwt"
)

func TestSubjectIDFromClaims_ForEmployee(t *testing.T) {
	t.Parallel()

	employeeID := uint(11)
	subjectID, ok, err := SubjectIDFromClaims(&jwt.Claims{
		IdentityType: string(IdentityEmployee),
		EmployeeID:   &employeeID,
	})

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, employeeID, subjectID)
}

func TestSubjectIDFromClaims_ForClient(t *testing.T) {
	t.Parallel()

	clientID := uint(22)
	subjectID, ok, err := SubjectIDFromClaims(&jwt.Claims{
		IdentityType: string(IdentityClient),
		ClientID:     &clientID,
	})

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, clientID, subjectID)
}

func TestSetSubjectIDOnClaims_ForClient(t *testing.T) {
	t.Parallel()

	claims := &jwt.Claims{IdentityType: string(IdentityClient)}
	err := SetSubjectIDOnClaims(claims, IdentityClient, 33)

	require.NoError(t, err)
	require.NotNil(t, claims.ClientID)
	require.Equal(t, uint(33), *claims.ClientID)
	require.Nil(t, claims.EmployeeID)
}

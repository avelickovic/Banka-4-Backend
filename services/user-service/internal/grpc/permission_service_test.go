package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/jwt"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/permission"
)

type fakePermissionProvider struct {
	captured *jwt.Claims
	result   []permission.Permission
	err      error
}

func (f *fakePermissionProvider) GetPermissions(_ context.Context, claims *jwt.Claims) ([]permission.Permission, error) {
	f.captured = claims
	return f.result, f.err
}

func TestPermissionService_GetPermissions_ForClient(t *testing.T) {
	t.Parallel()

	provider := &fakePermissionProvider{
		result: []permission.Permission{permission.Trading, permission.TradingMargin},
	}

	service := NewPermissionService(provider)

	resp, err := service.GetPermissions(context.Background(), &pb.GetPermissionsRequest{
		IdentityId:   42,
		IdentityType: "client",
		SubjectId:    7,
	})

	require.NoError(t, err)
	require.Equal(t, []string{string(permission.Trading), string(permission.TradingMargin)}, resp.GetPermissions())
	require.NotNil(t, provider.captured)
	require.Equal(t, uint(42), provider.captured.IdentityID)
	require.Equal(t, "client", provider.captured.IdentityType)
	require.NotNil(t, provider.captured.ClientID)
	require.Equal(t, uint(7), *provider.captured.ClientID)
	require.Nil(t, provider.captured.EmployeeID)
}

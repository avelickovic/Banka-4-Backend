package grpc

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/auth"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/jwt"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PermissionService struct {
	pb.UnimplementedPermissionServiceServer
	provider auth.PermissionProvider
}

func NewPermissionService(provider auth.PermissionProvider) *PermissionService {
	return &PermissionService{provider: provider}
}

func (s *PermissionService) GetPermissions(ctx context.Context, req *pb.GetPermissionsRequest) (*pb.GetPermissionsResponse, error) {
	claims := &jwt.Claims{
		IdentityID:   uint(req.GetIdentityId()),
		IdentityType: req.GetIdentityType(),
	}

	if req.GetSubjectId() > 0 {
		if err := auth.SetSubjectIDOnClaims(claims, auth.IdentityType(req.GetIdentityType()), uint(req.GetSubjectId())); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	permissions, err := s.provider.GetPermissions(ctx, claims)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(permissions))
	for _, p := range permissions {
		result = append(result, string(p))
	}

	return &pb.GetPermissionsResponse{Permissions: result}, nil
}

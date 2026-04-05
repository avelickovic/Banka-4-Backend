package client

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
)

type UserServiceClient interface {
	GetClientById(ctx context.Context, id uint64) (*pb.GetClientByIdResponse, error)
	GetEmployeeById(ctx context.Context, id uint64) (*pb.GetEmployeeByIdResponse, error)
	GetAllClients(ctx context.Context, page, pageSize int32, firstName, lastName string) (*pb.GetAllClientsResponse, error)
	GetAllActuaries(ctx context.Context, page, pageSize int32, firstName, lastName string) (*pb.GetAllActuariesResponse, error)
}

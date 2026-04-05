package grpc

import (
	"context"
	"fmt"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/client"
)

type UserServiceClient struct {
	stub pb.UserServiceClient
}

func NewUserServiceClient(conn *client.UserConn) *UserServiceClient {
	return &UserServiceClient{stub: pb.NewUserServiceClient(conn.ClientConn)}
}

func (c *UserServiceClient) GetClientById(ctx context.Context, id uint64) (*pb.GetClientByIdResponse, error) {
	resp, err := c.stub.GetClientById(ctx, &pb.GetClientByIdRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("user client GetClientById: %w", err)
	}
	return resp, nil
}

func (c *UserServiceClient) GetEmployeeById(ctx context.Context, id uint64) (*pb.GetEmployeeByIdResponse, error) {
	resp, err := c.stub.GetEmployeeById(ctx, &pb.GetEmployeeByIdRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("user client GetEmployeeById: %w", err)
	}
	return resp, nil
}

func (c *UserServiceClient) GetAllClients(ctx context.Context, page, pageSize int32, firstName, lastName string) (*pb.GetAllClientsResponse, error) {
	resp, err := c.stub.GetAllClients(ctx, &pb.GetAllClientsRequest{
		Page: page, PageSize: pageSize, FirstName: firstName, LastName: lastName,
	})
	if err != nil {
		return nil, fmt.Errorf("user client GetAllClients: %w", err)
	}
	return resp, nil
}

func (c *UserServiceClient) GetAllActuaries(ctx context.Context, page, pageSize int32, firstName, lastName string) (*pb.GetAllActuariesResponse, error) {
	resp, err := c.stub.GetAllActuaries(ctx, &pb.GetAllActuariesRequest{
		Page: page, PageSize: pageSize, FirstName: firstName, LastName: lastName,
	})
	if err != nil {
		return nil, fmt.Errorf("user client GetAllActuaries: %w", err)
	}
	return resp, nil
}

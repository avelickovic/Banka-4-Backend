package grpc

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/repository"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/service"
)

type BankingService struct {
	pb.UnimplementedBankingServiceServer
	accountRepo     repository.AccountRepository
	paymentService  *service.PaymentService
	exchangeService *service.ExchangeService
}

func NewBankingService(accountRepo repository.AccountRepository, paymentService *service.PaymentService, exchangeService *service.ExchangeService) *BankingService {
	return &BankingService{
		accountRepo:     accountRepo,
		paymentService:  paymentService,
		exchangeService: exchangeService,
	}
}

func (s *BankingService) GetAccountByNumber(ctx context.Context, req *pb.GetAccountByNumberRequest) (*pb.GetAccountByNumberResponse, error) {
	account, err := s.accountRepo.FindByAccountNumber(ctx, req.AccountNumber)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch account: %v", err)
	}
	if account == nil {
		return nil, status.Errorf(codes.NotFound, "account %s not found", req.AccountNumber)
	}
	return &pb.GetAccountByNumberResponse{
		AccountNumber:    account.AccountNumber,
		ClientId:         uint64(account.ClientID),
		AccountType:      string(account.AccountType),
		CurrencyCode:     string(account.Currency.Code),
		AvailableBalance: account.AvailableBalance,
	}, nil
}

func (s *BankingService) CreatePaymentWithoutVerification(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	payment, err := s.paymentService.CreatePaymentWithoutVerification(ctx, dto.CreatePaymentRequest{
		PayerAccountNumber:     req.PayerAccountNumber,
		RecipientAccountNumber: req.RecipientAccountNumber,
		RecipientName:          req.RecipientName,
		Amount:                 req.Amount,
		ReferenceNumber:        req.ReferenceNumber,
		PaymentCode:            req.PaymentCode,
		Purpose:                req.Purpose,
	})
	if err != nil {
		return nil, errors.MapGrpcToHttpError(err)
	}
	return &pb.CreatePaymentResponse{
		PaymentId:     uint64(payment.PaymentID),
		TransactionId: uint64(payment.TransactionID),
		Status:        string(payment.Transaction.Status),
	}, nil
}

func (s *BankingService) GetAccountsByClientID(ctx context.Context, req *pb.GetAccountsByClientIDRequest) (*pb.GetAccountsByClientIDResponse, error) {
	accounts, err := s.accountRepo.FindByClientID(ctx, uint(req.ClientId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch accounts: %v", err)
	}

	pbAccounts := make([]*pb.AccountInfo, 0, len(accounts))
	for _, acc := range accounts {
		pbAccounts = append(pbAccounts, &pb.AccountInfo{
			AccountNumber: acc.AccountNumber,
			CurrencyCode:  string(acc.Currency.Code),
		})
	}

	return &pb.GetAccountsByClientIDResponse{Accounts: pbAccounts}, nil
}

func (s *BankingService) ConvertCurrency(ctx context.Context, req *pb.ConvertCurrencyRequest) (*pb.ConvertCurrencyResponse, error) {
	if req.Amount <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "amount must be positive")
	}
	if req.FromCode == "" || req.ToCode == "" {
		return nil, status.Errorf(codes.InvalidArgument, "from_code and to_code are required")
	}

	converted, err := s.exchangeService.Convert(
		ctx,
		req.Amount,
		model.CurrencyCode(req.FromCode),
		model.CurrencyCode(req.ToCode),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "currency conversion failed: %v", err)
	}

	return &pb.ConvertCurrencyResponse{
		ConvertedAmount: converted,
	}, nil
}

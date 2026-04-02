package grpc

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apperrors "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/model"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/repository"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/service"
)

type BankingService struct {
	pb.UnimplementedBankingServiceServer
	accountRepo          repository.AccountRepository
	transactionRepo      repository.TransactionRepository
	transactionProcessor *service.TransactionProcessor
	exchangeService      service.CurrencyConverter
}

func NewBankingService(
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	transactionProcessor *service.TransactionProcessor,
	exchangeService service.CurrencyConverter,
) *BankingService {
	return &BankingService{
		accountRepo:          accountRepo,
		transactionRepo:      transactionRepo,
		transactionProcessor: transactionProcessor,
		exchangeService:      exchangeService,
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

func (s *BankingService) ExecuteTradeSettlement(ctx context.Context, req *pb.ExecuteTradeSettlementRequest) (*pb.ExecuteTradeSettlementResponse, error) {
	accountNumber := strings.TrimSpace(req.GetAccountNumber())
	if accountNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "account_number is required")
	}

	tradeCurrencyCode := strings.ToUpper(strings.TrimSpace(req.GetTradeCurrencyCode()))
	if tradeCurrencyCode == "" {
		return nil, status.Error(codes.InvalidArgument, "trade_currency_code is required")
	}

	direction := req.GetDirection()
	if direction == pb.TradeSettlementDirection_TRADE_SETTLEMENT_DIRECTION_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "direction is required")
	}

	amount := req.GetAmount()

	if amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be greater than zero")
	}

	customerAccount, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch account: %v", err)
	}
	if customerAccount == nil {
		return nil, status.Errorf(codes.NotFound, "account %s not found", accountNumber)
	}

	bankAccountNumber, err := resolveBankAccountNumber(tradeCurrencyCode)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	bankAccount, err := s.accountRepo.FindByAccountNumber(ctx, bankAccountNumber)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch bank settlement account: %v", err)
	}
	if bankAccount == nil {
		return nil, status.Errorf(codes.NotFound, "bank settlement account %s not found", bankAccountNumber)
	}

	sourceAccount := customerAccount
	destinationAccount := bankAccount
	amountIsSource := false

	switch direction {
	case pb.TradeSettlementDirection_TRADE_SETTLEMENT_DIRECTION_BUY:
	case pb.TradeSettlementDirection_TRADE_SETTLEMENT_DIRECTION_SELL:
		sourceAccount = bankAccount
		destinationAccount = customerAccount
		amountIsSource = true
	default:
		return nil, status.Error(codes.InvalidArgument, "direction must be BUY or SELL")
	}

	sourceAmount := amount
	destinationAmount := amount
	if sourceAccount.Currency.Code != destinationAccount.Currency.Code {
		if amountIsSource {
			destinationAmount, err = s.exchangeService.Convert(ctx, amount, sourceAccount.Currency.Code, destinationAccount.Currency.Code)
		} else {
			sourceAmount, err = s.exchangeService.Convert(ctx, amount, destinationAccount.Currency.Code, sourceAccount.Currency.Code)
		}
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert currencies: %v", err)
		}
	}

	transaction := &model.Transaction{
		PayerAccountNumber:     sourceAccount.AccountNumber,
		RecipientAccountNumber: destinationAccount.AccountNumber,
		StartAmount:            sourceAmount,
		StartCurrencyCode:      sourceAccount.Currency.Code,
		EndAmount:              destinationAmount,
		EndCurrencyCode:        destinationAccount.Currency.Code,
		Status:                 model.TransactionProcessing,
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create transaction: %v", err)
	}

	if err := s.transactionProcessor.ProcessTradeSettlement(ctx, transaction.TransactionID); err != nil {
		return nil, mapTradeSettlementError(err)
	}

	return &pb.ExecuteTradeSettlementResponse{
		TransactionId:           uint64(transaction.TransactionID),
		SourceAmount:            sourceAmount,
		SourceCurrencyCode:      string(sourceAccount.Currency.Code),
		DestinationAmount:       destinationAmount,
		DestinationCurrencyCode: string(destinationAccount.Currency.Code),
	}, nil
}

func resolveBankAccountNumber(currencyCode string) (string, error) {
	accountNumber, ok := service.BankAccounts[model.CurrencyCode(currencyCode)]
	if !ok {
		return "", fmt.Errorf("unsupported trade currency: %s", currencyCode)
	}

	return accountNumber, nil
}

func mapTradeSettlementError(err error) error {
	var appErr *apperrors.AppError
	if stderrors.As(err, &appErr) {
		switch appErr.Code {
		case 400:
			return status.Error(codes.FailedPrecondition, appErr.Message)
		case 404:
			return status.Error(codes.NotFound, appErr.Message)
		default:
			return status.Error(codes.Internal, appErr.Message)
		}
	}

	return status.Error(codes.Internal, err.Error())
}

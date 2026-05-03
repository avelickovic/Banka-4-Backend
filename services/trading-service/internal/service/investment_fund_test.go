package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/auth"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/stretchr/testify/require"
)

// ── Fake Fund Repo (extended for GetFundDetail) ───────────────────────────────

type fakeFundRepo struct {
	findByIDResult   *model.InvestmentFund
	findByIDErr      error
	findByNameResult *model.InvestmentFund
	findByNameErr    error
	createErr        error
	created          *model.InvestmentFund

	// new fields for GetFundDetail
	findHoldingsResult          []model.AssetOwnership
	findHoldingsErr             error
	getPerformanceHistoryResult []model.FundPerformance
	getPerformanceHistoryErr    error

	findAllResult       []model.InvestmentFund
	findAllTotal        int64
	findAllErr          error
	findByManagerResult []model.InvestmentFund
	findByManagerErr    error
}

func (f *fakeFundRepo) FindByName(ctx context.Context, name string) (*model.InvestmentFund, error) {
	return f.findByNameResult, f.findByNameErr
}
func (f *fakeFundRepo) FindByID(ctx context.Context, id uint) (*model.InvestmentFund, error) {
	return f.findByIDResult, f.findByIDErr
}
func (f *fakeFundRepo) FindByAccountNumber(ctx context.Context, accountNumber string) (*model.InvestmentFund, error) {
	return nil, nil
}

func (f *fakeFundRepo) GetAllInvestmentFunds(ctx context.Context) ([]model.InvestmentFund, error) {
	return f.findAllResult, f.findAllErr
}

func (f *fakeFundRepo) FindAll(ctx context.Context, name, sortBy, sortDir string, page, pageSize int) ([]model.InvestmentFund, int64, error) {
	return f.findAllResult, f.findAllTotal, f.findAllErr
}

func (f *fakeFundRepo) FindByManagerID(ctx context.Context, managerID uint) ([]model.InvestmentFund, error) {
	return f.findByManagerResult, f.findByManagerErr
}

func (f *fakeFundRepo) Create(ctx context.Context, fund *model.InvestmentFund) error {
	if f.createErr != nil {
		return f.createErr
	}
	fund.FundID = 1
	f.created = fund
	return nil
}
func (f *fakeFundRepo) FindHoldings(ctx context.Context, fundID uint) ([]model.AssetOwnership, error) {
	return f.findHoldingsResult, f.findHoldingsErr
}

func (f *fakeFundRepo) GetPerformanceHistory(ctx context.Context, fundID uint, limit int) ([]model.FundPerformance, error) {
	return f.getPerformanceHistoryResult, f.getPerformanceHistoryErr
}
func (f *fakeFundRepo) SavePerformanceSnapshot(ctx context.Context, perf *model.FundPerformance) error {
	return nil
}

// ── Fake Position / Investment Repos (unchanged) ─────────────────────────────

type fakePositionRepo struct {
	findResult      *model.ClientFundPosition
	findErr         error
	findByClientRes []model.ClientFundPosition
	findByClientErr error
	findByFundRes   []model.ClientFundPosition
	findByFundErr   error
	upsertErr       error
}

func (f *fakePositionRepo) FindByClientAndFund(ctx context.Context, clientID uint, ownerType model.OwnerType, fundID uint) (*model.ClientFundPosition, error) {
	return f.findResult, f.findErr
}

func (f *fakePositionRepo) FindByClient(ctx context.Context, clientID uint, ownerType model.OwnerType) ([]model.ClientFundPosition, error) {
	return f.findByClientRes, f.findByClientErr
}

func (f *fakePositionRepo) FindByFund(ctx context.Context, fundID uint) ([]model.ClientFundPosition, error) {
	return f.findByFundRes, f.findByFundErr
}
func (f *fakePositionRepo) Upsert(ctx context.Context, position *model.ClientFundPosition) error {
	return f.upsertErr
}

type fakeInvestmentRepo struct {
	createErr error
}

func (f *fakeInvestmentRepo) Create(ctx context.Context, investment *model.ClientFundInvestment) error {
	return f.createErr
}
func (f *fakeInvestmentRepo) FindByClientAndFund(ctx context.Context, clientID uint, ownerType model.OwnerType, fundID uint) ([]model.ClientFundInvestment, error) {
	return nil, nil
}
func (f *fakeInvestmentRepo) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

// ── Fake Banking Client (unchanged) ─────────────────────────────────────────

type fakeFundBankingClient struct {
	createdAccountNumber string
	createFundAccountErr error
	getAccountResult     *pb.GetAccountByNumberResponse
	tradeSettlementErr   error
	convertCurrencyFunc  func(amount float64, from, to string) (float64, error)
}

func (f *fakeFundBankingClient) GetAccountByNumber(_ context.Context, _ string) (*pb.GetAccountByNumberResponse, error) {
	return f.getAccountResult, nil
}
func (f *fakeFundBankingClient) HasActiveLoan(_ context.Context, _ uint64) (*pb.HasActiveLoanResponse, error) {
	return nil, nil
}
func (f *fakeFundBankingClient) CreatePaymentWithoutVerification(_ context.Context, _ *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	return nil, nil
}
func (f *fakeFundBankingClient) GetAccountsByClientID(_ context.Context, _ uint64) (*pb.GetAccountsByClientIDResponse, error) {
	return nil, nil
}
func (f *fakeFundBankingClient) ConvertCurrency(_ context.Context, amount float64, from, to string) (float64, error) {
	if f.convertCurrencyFunc != nil {
		return f.convertCurrencyFunc(amount, from, to)
	}
	return amount, nil
}
func (f *fakeFundBankingClient) ExecuteTradeSettlement(_ context.Context, _, _ string, _ pb.TradeSettlementDirection, _ float64) (*pb.ExecuteTradeSettlementResponse, error) {
	if f.tradeSettlementErr != nil {
		return nil, f.tradeSettlementErr
	}
	return &pb.ExecuteTradeSettlementResponse{}, nil
}
func (f *fakeFundBankingClient) GetAccountCurrency(_ context.Context, _ string) (string, error) {
	return "RSD", nil
}
func (f *fakeFundBankingClient) CreateFundAccount(_ context.Context, _ string, _ uint64) (string, error) {
	return f.createdAccountNumber, f.createFundAccountErr
}

// ── Fake User Client
type fakeUserClient struct {
	getEmployeeByIDFunc func(ctx context.Context, id uint64) (*pb.GetEmployeeByIdResponse, error)
}

type fakeFundUserClient struct {
	// configurable responses
	getClientByIdResp *pb.GetClientByIdResponse
	getClientByIdErr  error

	getEmployeeByIdResp *pb.GetEmployeeByIdResponse
	getEmployeeByIdErr  error

	getAllClientsResp *pb.GetAllClientsResponse
	getAllClientsErr  error

	getAllActuariesResp *pb.GetAllActuariesResponse
	getAllActuariesErr  error

	getIdentityByUserIdResp *pb.GetIdentityByUserIdResponse
	getIdentityByUserIdErr  error
}

func (f *fakeFundUserClient) GetClientById(_ context.Context, _ uint64) (*pb.GetClientByIdResponse, error) {
	return f.getClientByIdResp, f.getClientByIdErr
}

func (f *fakeFundUserClient) GetClientByIdentityId(_ context.Context, _ uint64) (*pb.GetClientByIdResponse, error) {
	return f.getClientByIdResp, f.getClientByIdErr
}

func (f *fakeFundUserClient) GetEmployeeById(_ context.Context, _ uint64) (*pb.GetEmployeeByIdResponse, error) {
	return f.getEmployeeByIdResp, f.getEmployeeByIdErr
}

func (f *fakeFundUserClient) GetEmployeeByIdentityId(_ context.Context, _ uint64) (*pb.GetEmployeeByIdResponse, error) {
	return f.getEmployeeByIdResp, f.getEmployeeByIdErr
}

func (f *fakeFundUserClient) GetAllClients(_ context.Context, _, _ int32, _, _ string) (*pb.GetAllClientsResponse, error) {
	return f.getAllClientsResp, f.getAllClientsErr
}

func (f *fakeFundUserClient) GetAllActuaries(_ context.Context, _, _ int32, _, _ string) (*pb.GetAllActuariesResponse, error) {
	return f.getAllActuariesResp, f.getAllActuariesErr
}

func (f *fakeFundUserClient) GetIdentityByUserId(_ context.Context, _ uint64, _ string) (*pb.GetIdentityByUserIdResponse, error) {
	return f.getIdentityByUserIdResp, f.getIdentityByUserIdErr
}

// ── Helpers ───────────────────────────────────────────────────────

func fundSupervisorCtx() context.Context {
	employeeID := uint(25)
	return auth.SetAuthOnContext(context.Background(), &auth.AuthContext{
		IdentityID:   200,
		IdentityType: auth.IdentityEmployee,
		EmployeeID:   &employeeID,
	})
}

func (f *fakeUserClient) GetClientById(ctx context.Context, id uint64) (*pb.GetClientByIdResponse, error) {
	return nil, nil
}
func (f *fakeUserClient) GetClientByIdentityId(ctx context.Context, identityId uint64) (*pb.GetClientByIdResponse, error) {
	return nil, nil
}
func (f *fakeUserClient) GetEmployeeById(ctx context.Context, id uint64) (*pb.GetEmployeeByIdResponse, error) {
	if f.getEmployeeByIDFunc != nil {
		return f.getEmployeeByIDFunc(ctx, id)
	}
	return &pb.GetEmployeeByIdResponse{Id: id, FullName: fmt.Sprintf("Manager %d", id)}, nil
}
func (f *fakeUserClient) GetEmployeeByIdentityId(ctx context.Context, identityId uint64) (*pb.GetEmployeeByIdResponse, error) {
	return nil, nil
}
func (f *fakeUserClient) GetAllClients(ctx context.Context, page, pageSize int32, firstName, lastName string) (*pb.GetAllClientsResponse, error) {
	return nil, nil
}
func (f *fakeUserClient) GetAllActuaries(ctx context.Context, page, pageSize int32, firstName, lastName string) (*pb.GetAllActuariesResponse, error) {
	return nil, nil
}
func (f *fakeUserClient) GetIdentityByUserId(ctx context.Context, userID uint64, userType string) (*pb.GetIdentityByUserIdResponse, error) {
	return nil, nil
}

// ── Helper for creating service with listingRepo ───────────────────────────

func newTestFundServiceWithListing(fundRepo *fakeFundRepo, listingRepo *fakeListingRepo, bankingClient *fakeFundBankingClient, userClient *fakeUserClient) *InvestmentFundService {
	exchange := defaultExchange()
	svc := NewInvestmentFundService(fundRepo, &fakePositionRepo{}, listingRepo, &fakeInvestmentRepo{}, &fakeAssetOwnershipRepo{}, &fakeExchangeRepo{exchange: exchange}, &fakeStockRepo{},
		&fakeOptionRepo{},
		&fakeFuturesRepo{},
		&fakeForexRepo{}, bankingClient, userClient)
	svc.listingRepo = listingRepo // inject listingRepo
	return svc
}

// ── Tests: GetFundDetail ───────────────────────────────────────────────────

func TestGetFundDetail_Success(t *testing.T) {
	fund := &model.InvestmentFund{
		FundID:              1,
		Name:                "Test Fund",
		Description:         "A test fund",
		MinimumContribution: 500,
		ManagerID:           10,
		AccountNumber:       "ACC123",
		Positions: []model.ClientFundPosition{
			{TotalInvestedAmount: 1500.0},
		},
	}
	bankingClient := &fakeFundBankingClient{
		getAccountResult: &pb.GetAccountByNumberResponse{AvailableBalance: 2000},
	}
	userClient := &fakeUserClient{}
	fundRepo := &fakeFundRepo{
		findByIDResult: fund,
		findHoldingsResult: []model.AssetOwnership{
			{
				AssetID:   100,
				Amount:    10,
				Asset:     model.Asset{Ticker: "AAPL"}, // not pointer
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			{
				AssetID:   101,
				Amount:    5,
				Asset:     model.Asset{Ticker: "GOOGL"}, // not pointer
				UpdatedAt: time.Now().Add(-48 * time.Hour),
			},
		},
		getPerformanceHistoryResult: []model.FundPerformance{
			{Date: time.Now().AddDate(0, -1, 0), FundValue: 1800},
			{Date: time.Now().AddDate(0, -2, 0), FundValue: 1700},
		},
	}
	listingRepo := &fakeListingRepo{
		byAssetIDs: []model.Listing{
			{AssetID: 100, Price: 120, MaintenanceMargin: 10, ListingID: 1000, Exchange: &model.Exchange{MicCode: "XSIM"}},
			{AssetID: 101, Price: 110, MaintenanceMargin: 8, ListingID: 1001, Exchange: &model.Exchange{MicCode: "XSIM"}},
		},
		dailyPriceInfo: &model.ListingDailyPriceInfo{Change: 2.5, Volume: 1000},
	}
	svc := newTestFundServiceWithListing(fundRepo, listingRepo, bankingClient, userClient)

	resp, err := svc.GetFundDetail(context.Background(), 1)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "Test Fund", resp.Name)
	require.Equal(t, "A test fund", resp.Description)
	require.Equal(t, 500.0, resp.MinInvestment)
	require.Equal(t, "Manager 10", resp.Manager)
	require.Equal(t, 2000.0, resp.LiquidAssets)

	// fundValue = (10*120)+(5*110) = 1200+550=1750
	// profit = 1750 - 1500 = 250
	require.Equal(t, 3750.0, resp.FundValue)
	require.Equal(t, 2250.0, resp.Profit)

	require.Len(t, resp.Holdings, 2)
	require.Equal(t, "AAPL", resp.Holdings[0].Ticker)
	require.Equal(t, 120.0, resp.Holdings[0].Price)
	require.Equal(t, 2.5, resp.Holdings[0].Change)
	require.Equal(t, uint64(1000), resp.Holdings[0].Volume)
	require.Equal(t, 10.0, resp.Holdings[0].InitialMarginCost)

	require.Len(t, resp.PerformanceHistory, 2)
}

func newTestFundService(fundRepo *fakeFundRepo, ownershipRepo *fakeAssetOwnershipRepo, listingRepo *fakeListingRepo, bankingClient *fakeFundBankingClient, userClient *fakeFundUserClient) *InvestmentFundService {
	exchange := defaultExchange()
	return NewInvestmentFundService(
		fundRepo,
		&fakePositionRepo{}, listingRepo,
		&fakeInvestmentRepo{},
		ownershipRepo,
		&fakeExchangeRepo{exchange: exchange},
		&fakeStockRepo{},
		&fakeOptionRepo{},
		&fakeFuturesRepo{},
		&fakeForexRepo{},
		bankingClient,
		userClient,
	)
}

// ── CreateFund tests ──────────────────────────────────────────────

func TestCreateFund_Success(t *testing.T) {
	fundRepo := &fakeFundRepo{}
	bankingClient := &fakeFundBankingClient{createdAccountNumber: "444000112345678901"}
	svc := newTestFundService(fundRepo, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, bankingClient, &fakeFundUserClient{})

	resp, err := svc.CreateFund(fundSupervisorCtx(), validFundRequest())

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "Alpha Growth Fund", resp.Name)
	require.Equal(t, "444000112345678901", resp.AccountNumber)
	require.Equal(t, uint(25), resp.ManagerID)
	require.Equal(t, 1000.00, resp.MinimumContribution)
	require.WithinDuration(t, time.Now(), resp.CreatedAt, 5*time.Second)
}

func TestGetFundDetail_NotFound(t *testing.T) {
	fundRepo := &fakeFundRepo{findByIDResult: nil}
	svc := newTestFundServiceWithListing(fundRepo, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeUserClient{})
	_, err := svc.GetFundDetail(context.Background(), 99)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestCreateFund_Unauthenticated(t *testing.T) {
	svc := newTestFundService(&fakeFundRepo{}, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	_, err := svc.CreateFund(context.Background(), validFundRequest())

	require.Error(t, err)
	require.Contains(t, err.Error(), "not authenticated")
}

func TestGetFundDetail_RepoFindByIDError(t *testing.T) {
	fundRepo := &fakeFundRepo{findByIDErr: errors.New("db error")}
	svc := newTestFundServiceWithListing(fundRepo, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeUserClient{})
	_, err := svc.GetFundDetail(context.Background(), 1)

	require.Error(t, err)
}

func TestCreateFund_NotEmployee(t *testing.T) {
	svc := newTestFundService(&fakeFundRepo{}, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeFundUserClient{})

	_, err := svc.CreateFund(fundClientCtx(), validFundRequest())

	require.Error(t, err)
}

func TestGetFundDetail_HoldingsError(t *testing.T) {
	fund := &model.InvestmentFund{FundID: 1, AccountNumber: "ACC"}
	fundRepo := &fakeFundRepo{
		findByIDResult:  fund,
		findHoldingsErr: errors.New("holdings error"),
	}
	svc := newTestFundServiceWithListing(fundRepo, &fakeListingRepo{}, &fakeFundBankingClient{}, &fakeUserClient{})
	_, err := svc.GetFundDetail(context.Background(), 1)
	require.Error(t, err)
}

func TestCreateFund_BankingClientError(t *testing.T) {
	fundRepo := &fakeFundRepo{}
	bankingClient := &fakeFundBankingClient{
		createFundAccountErr: fmt.Errorf("banking service unavailable"),
	}
	svc := newTestFundService(fundRepo, &fakeAssetOwnershipRepo{}, &fakeListingRepo{}, bankingClient, &fakeFundUserClient{})

	_, err := svc.CreateFund(fundSupervisorCtx(), validFundRequest())

	require.Error(t, err)
}

func TestGetFundDetail_EmptyHoldings(t *testing.T) {
	fund := &model.InvestmentFund{FundID: 1, AccountNumber: "ACC", MinimumContribution: 100}
	fundRepo := &fakeFundRepo{
		findByIDResult:     fund,
		findHoldingsResult: []model.AssetOwnership{},
		getPerformanceHistoryResult: []model.FundPerformance{
			{Date: time.Now(), FundValue: 5000},
		},
	}
	listingRepo := &fakeListingRepo{}
	bankingClient := &fakeFundBankingClient{
		getAccountResult: &pb.GetAccountByNumberResponse{AvailableBalance: 5000},
	}
	userClient := &fakeUserClient{}
	svc := newTestFundServiceWithListing(fundRepo, listingRepo, bankingClient, userClient)

	resp, err := svc.GetFundDetail(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, 5000.0, resp.FundValue)
	require.Equal(t, 5000.0, resp.Profit) // profit = 5000 - 0 = 5000
	require.Empty(t, resp.Holdings)
	require.NotEmpty(t, resp.PerformanceHistory)
}

func validFundRequest() dto.CreateFundRequest {
	return dto.CreateFundRequest{
		Name:                "Alpha Growth Fund",
		Description:         "Fund focused on the IT sector.",
		MinimumContribution: 1000.00,
	}
}

func fundClientCtx() context.Context {
	clientID := uint(99)
	return auth.SetAuthOnContext(context.Background(), &auth.AuthContext{
		IdentityID:   1,
		IdentityType: auth.IdentityClient,
		ClientID:     &clientID,
	})
}

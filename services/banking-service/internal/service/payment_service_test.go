package service

import (
	"banking-service/internal/dto"
	"banking-service/internal/model"
	"banking-service/internal/repository"
	"common/pkg/auth"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ── Fake Payment Repo ──────────────────────────────────────────────

type fakePaymentRepo struct {
	createErr error
	getErr    error
	payment   *model.Payment
}

func (f *fakePaymentRepo) Create(ctx context.Context, p *model.Payment) error {
	if f.createErr != nil {
		return f.createErr
	}
	p.PaymentID = 1
	f.payment = p
	return nil
}

func (f *fakePaymentRepo) GetByID(ctx context.Context, id uint) (*model.Payment, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.payment, nil
}

func (f *fakePaymentRepo) Update(ctx context.Context, p *model.Payment) error {
	f.payment = p
	return nil
}

type fakeTransactionRepo struct {
	createErr     error
	getErr        error
	updateErr     error
	transaction   *model.Transaction
	returnedTx    *model.Transaction
	returnedTxErr error
}

func (f *fakeTransactionRepo) Create(_ context.Context, t *model.Transaction) error {
	if f.createErr != nil {
		return f.createErr
	}
	t.TransactionID = 1 // simulate ID assignment
	f.transaction = t
	return nil
}

func (f *fakeTransactionRepo) Update(ctx context.Context, t *model.Transaction) error {
	f.transaction = t
	return nil
}

func (f *fakeTransactionRepo) GetByID(_ context.Context, _ uint) (*model.Transaction, error) {
	// return preset transaction or error
	return f.returnedTx, f.returnedTxErr
}

// ── Fake Payment Account Repo ──────────────────────────────────────

type fakePaymentAccountRepo struct {
	accounts map[string]*model.Account
	findErr  error
}

func newFakePaymentAccountRepo(accounts ...*model.Account) *fakePaymentAccountRepo {
	m := make(map[string]*model.Account)
	for _, a := range accounts {
		m[a.AccountNumber] = a
	}
	return &fakePaymentAccountRepo{accounts: m}
}

func (f *fakePaymentAccountRepo) Create(ctx context.Context, account *model.Account) error {
	return nil
}

func (f *fakePaymentAccountRepo) AccountNumberExists(ctx context.Context, accountNumber string) (bool, error) {
	_, exists := f.accounts[accountNumber]
	return exists, nil
}

func (f *fakePaymentAccountRepo) FindByAccountNumber(ctx context.Context, accountNumber string) (*model.Account, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	acc, exists := f.accounts[accountNumber]
	if !exists {
		return nil, errors.New("account not found")
	}
	return acc, nil
}

func (f *fakePaymentAccountRepo) UpdateBalance(ctx context.Context, account *model.Account) error {
	f.accounts[account.AccountNumber] = account
	return nil
}

// ── Fake Exchange Service ──────────────────────────────────────────

type fakeExchangeService struct {
	rate float64
	err  error
}

func (f *fakeExchangeService) Convert(ctx context.Context, amount float64, from, to model.CurrencyCode) (float64, error) {
	if f.err != nil {
		return 0, f.err
	}
	return amount * f.rate, nil
}

type fakeMobileSecretClient struct {
	secret string
	err    error
}

func (f *fakeMobileSecretClient) GetMobileSecret(_ context.Context, _ string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	if f.secret != "" {
		return f.secret, nil
	}
	return "JBSWY3DPEHPK3PXP", nil
}

type fakeVerifyTransactionProcessor struct {
	err       error
	processed []uint
}

func (f *fakeVerifyTransactionProcessor) Process(_ context.Context, transactionID uint) error {
	f.processed = append(f.processed, transactionID)
	return f.err
}

// ── Constructor ────────────────────────────────────────────────────────

func newTestPaymentService(
	paymentRepo repository.PaymentRepository,
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	exchangeSvc CurrencyConverter,
) *PaymentService {
	return &PaymentService{
		paymentRepo:          paymentRepo,
		transactionRepo:      transactionRepo,
		accountRepo:          accountRepo,
		exchangeService:      exchangeSvc,
		mobileSecretClient:   &fakeMobileSecretClient{},
		transactionProcessor: &fakeVerifyTransactionProcessor{},
		now:                  time.Now,
	}
}

func uintPtr(v uint) *uint {
	return &v
}

// ── Tests ──────────────────────────────────────────────────────────

func TestCreatePayment_Success(t *testing.T) {
	payerAccount := &model.Account{
		AccountNumber:    "87654321",
		ClientID:         1,
		Balance:          10000,
		AvailableBalance: 10000,
		DailyLimit:       250000,
		MonthlyLimit:     1000000,
		DailySpending:    0,
		MonthlySpending:  0,
		Currency:         model.Currency{Code: model.RSD},
	}
	recipientAccount := &model.Account{
		AccountNumber: "12345678",
		ClientID:      2,
		Currency:      model.Currency{Code: model.RSD},
	}

	svc := newTestPaymentService(
		&fakePaymentRepo{},
		&fakeTransactionRepo{},
		newFakePaymentAccountRepo(payerAccount, recipientAccount),
		&fakeExchangeService{rate: 1.0},
	)

	req := dto.CreatePaymentRequest{
		RecipientName:          "John Doe",
		RecipientAccountNumber: "12345678",
		Amount:                 100,
		PayerAccountNumber:     "87654321",
	}

	payment, err := svc.CreatePayment(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "John Doe", payment.RecipientName)
}

func TestCreatePayment_InsufficientFunds(t *testing.T) {
	payerAccount := &model.Account{
		AccountNumber:    "87654321",
		ClientID:         1,
		Balance:          50,
		AvailableBalance: 50,
		DailyLimit:       250000,
		MonthlyLimit:     1000000,
		Currency:         model.Currency{Code: model.RSD},
	}
	recipientAccount := &model.Account{
		AccountNumber: "12345678",
		ClientID:      2,
		Currency:      model.Currency{Code: model.RSD},
	}

	svc := newTestPaymentService(
		&fakePaymentRepo{},
		&fakeTransactionRepo{},
		newFakePaymentAccountRepo(payerAccount, recipientAccount),
		&fakeExchangeService{rate: 1.0},
	)

	req := dto.CreatePaymentRequest{
		RecipientName:          "John Doe",
		RecipientAccountNumber: "12345678",
		Amount:                 100,
		PayerAccountNumber:     "87654321",
	}

	payment, err := svc.CreatePayment(context.Background(), req)
	require.Nil(t, payment)
	require.Error(t, err)
	require.Contains(t, err.Error(), "insufficient funds")
}

func TestCreatePayment_DailyLimitExceeded(t *testing.T) {
	payerAccount := &model.Account{
		AccountNumber:    "87654321",
		ClientID:         1,
		Balance:          500000,
		AvailableBalance: 500000,
		DailyLimit:       1000,
		MonthlyLimit:     1000000,
		DailySpending:    900,
		MonthlySpending:  0,
		Currency:         model.Currency{Code: model.RSD},
	}
	recipientAccount := &model.Account{
		AccountNumber: "12345678",
		ClientID:      2,
		Currency:      model.Currency{Code: model.RSD},
	}

	svc := newTestPaymentService(
		&fakePaymentRepo{},
		&fakeTransactionRepo{},
		newFakePaymentAccountRepo(payerAccount, recipientAccount),
		&fakeExchangeService{rate: 1.0},
	)

	req := dto.CreatePaymentRequest{
		RecipientName:          "John Doe",
		RecipientAccountNumber: "12345678",
		Amount:                 200,
		PayerAccountNumber:     "87654321",
	}

	payment, err := svc.CreatePayment(context.Background(), req)
	require.Nil(t, payment)
	require.Error(t, err)
	require.Contains(t, err.Error(), "daily limit exceeded")
}

func TestCreatePayment_MonthlyLimitExceeded(t *testing.T) {
	payerAccount := &model.Account{
		AccountNumber:    "87654321",
		ClientID:         1,
		Balance:          500000,
		AvailableBalance: 500000,
		DailyLimit:       1000000,
		MonthlyLimit:     1000,
		DailySpending:    0,
		MonthlySpending:  900,
		Currency:         model.Currency{Code: model.RSD},
	}
	recipientAccount := &model.Account{
		AccountNumber: "12345678",
		ClientID:      2,
		Currency:      model.Currency{Code: model.RSD},
	}

	svc := newTestPaymentService(
		&fakePaymentRepo{},
		&fakeTransactionRepo{},
		newFakePaymentAccountRepo(payerAccount, recipientAccount),
		&fakeExchangeService{rate: 1.0},
	)

	req := dto.CreatePaymentRequest{
		RecipientName:          "John Doe",
		RecipientAccountNumber: "12345678",
		Amount:                 200,
		PayerAccountNumber:     "87654321",
	}

	payment, err := svc.CreatePayment(context.Background(), req)
	require.Nil(t, payment)
	require.Error(t, err)
	require.Contains(t, err.Error(), "monthly limit exceeded")
}

func TestCreatePayment_RecipientNotFound(t *testing.T) {
	payerAccount := &model.Account{
		AccountNumber:    "87654321",
		ClientID:         1,
		Balance:          10000,
		AvailableBalance: 10000,
		DailyLimit:       250000,
		MonthlyLimit:     1000000,
		Currency:         model.Currency{Code: model.RSD},
	}

	svc := newTestPaymentService(
		&fakePaymentRepo{},
		&fakeTransactionRepo{},
		newFakePaymentAccountRepo(payerAccount),
		&fakeExchangeService{rate: 1.0},
	)

	req := dto.CreatePaymentRequest{
		RecipientName:          "John Doe",
		RecipientAccountNumber: "99999999",
		Amount:                 100,
		PayerAccountNumber:     "87654321",
	}

	payment, err := svc.CreatePayment(context.Background(), req)
	require.Nil(t, payment)
	require.Error(t, err)
}

func TestCreatePayment_TransactionRepoError(t *testing.T) {
	payerAccount := &model.Account{
		AccountNumber:    "87654321",
		ClientID:         1,
		Balance:          10000,
		AvailableBalance: 10000,
		DailyLimit:       250000,
		MonthlyLimit:     1000000,
		Currency:         model.Currency{Code: model.RSD},
	}
	recipientAccount := &model.Account{
		AccountNumber: "12345678",
		ClientID:      2,
		Currency:      model.Currency{Code: model.RSD},
	}

	svc := newTestPaymentService(
		&fakePaymentRepo{},
		&fakeTransactionRepo{createErr: errors.New("db error")},
		newFakePaymentAccountRepo(payerAccount, recipientAccount),
		&fakeExchangeService{rate: 1.0},
	)

	req := dto.CreatePaymentRequest{
		RecipientName:          "John Doe",
		RecipientAccountNumber: "12345678",
		Amount:                 100,
		PayerAccountNumber:     "87654321",
	}

	payment, err := svc.CreatePayment(context.Background(), req)
	require.Nil(t, payment)
	require.Error(t, err)
}

func TestVerifyPayment_InvalidCode(t *testing.T) {
	paymentRepo := &fakePaymentRepo{
		payment: &model.Payment{
			PaymentID: 1,
			Transaction: model.Transaction{
				TransactionID:      42,
				PayerAccountNumber: "87654321",
				Status:             model.TransactionProcessing,
			},
		},
	}

	processor := &fakeVerifyTransactionProcessor{}
	authCtx := auth.SetAuthOnContext(context.Background(), &auth.AuthContext{ClientID: uintPtr(11)})
	svc := &PaymentService{
		paymentRepo:          paymentRepo,
		accountRepo:          newFakePaymentAccountRepo(&model.Account{AccountNumber: "87654321", ClientID: 11}),
		mobileSecretClient:   &fakeMobileSecretClient{secret: "JBSWY3DPEHPK3PXP"},
		transactionProcessor: processor,
		now:                  func() time.Time { return time.Unix(59, 0) },
	}

	payment, err := svc.VerifyPayment(authCtx, 1, "000000", "Bearer token")
	require.Nil(t, payment)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid verification code")
	require.Empty(t, processor.processed)
}

func TestVerifyPayment_Success(t *testing.T) {
	paymentRepo := &fakePaymentRepo{
		payment: &model.Payment{
			PaymentID: 1,
			Transaction: model.Transaction{
				TransactionID:      7,
				PayerAccountNumber: "87654321",
				Status:             model.TransactionProcessing,
			},
		},
	}

	processor := &fakeVerifyTransactionProcessor{}
	authCtx := auth.SetAuthOnContext(context.Background(), &auth.AuthContext{ClientID: uintPtr(5)})
	svc := &PaymentService{
		paymentRepo:          paymentRepo,
		accountRepo:          newFakePaymentAccountRepo(&model.Account{AccountNumber: "87654321", ClientID: 5}),
		mobileSecretClient:   &fakeMobileSecretClient{secret: "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"},
		transactionProcessor: processor,
		now:                  func() time.Time { return time.Unix(59, 0) },
	}

	payment, err := svc.VerifyPayment(authCtx, 1, "287082", "Bearer token")
	require.NoError(t, err)
	require.NotNil(t, payment)
	require.Equal(t, []uint{7}, processor.processed)
}

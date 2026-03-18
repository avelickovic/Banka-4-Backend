package service

import (
	"banking-service/internal/dto"
	"banking-service/internal/model"
	"banking-service/internal/repository"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// ── Fake Repo ────────────────────────────────────────────────────────

type fakePaymentRepo struct {
	createErr  error
	getErr     error
	payment    *model.Payment
	payments   []model.Payment
	findAccErr error
	total      int64
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
	createErr       error
	getErr          error
	updateErr       error
	transaction     *model.Transaction
	returnedTx      *model.Transaction
	returnedTxErr   error
}

func (f *fakeTransactionRepo) Create(_ context.Context, t *model.Transaction) error {
	if f.createErr != nil {
		return f.createErr
	}
	t.TransactionID = 1 // simulate ID assignment
	f.transaction = t
	return nil
}

func (f *fakePaymentRepo) FindByAccount(_ context.Context, _ string, _ *dto.PaymentFilters) ([]model.Payment, int64, error) {
	return f.payments, f.total, f.findAccErr
}

func (f *fakeTransactionRepo) GetByID(_ context.Context, _ uint) (*model.Transaction, error) {
	// return preset transaction or error
	return f.returnedTx, f.returnedTxErr
}

func (f *fakeTransactionRepo) Update(_ context.Context, _ *model.Transaction) error {
	return f.updateErr
}

// ── Constructor ────────────────────────────────────────────────────────

func newPaymentService(paymentRepo repository.PaymentRepository, transactionRepo repository.TransactionRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo, transactionRepo: transactionRepo}
}

// ── Tests ──────────────────────────────────────────────────────────────

func TestCreatePayment(t *testing.T) {
	paymentRepo := &fakePaymentRepo{}
	transactionRepo := &fakeTransactionRepo{}
	svc := newPaymentService(paymentRepo, transactionRepo)

	req := dto.CreatePaymentRequest{
		RecipientName:          "John Doe",
		RecipientAccountNumber: "12345678",
		Amount:                 100,
		PayerAccountNumber:     "87654321",
		CurrencyCode:           model.CurrencyCode("RSD"),
	}

	payment, err := svc.CreatePayment(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "John Doe", payment.RecipientName)
}

func TestCreatePayment_Error(t *testing.T) {
	paymentRepo := &fakePaymentRepo{createErr: errors.New("db error")}
	transactionRepo := &fakeTransactionRepo{}
	svc := newPaymentService(paymentRepo, transactionRepo)

	req := dto.CreatePaymentRequest{
		RecipientName:          "John Doe",
		RecipientAccountNumber: "12345678",
		Amount:                 100,
		PayerAccountNumber:     "87654321",
		CurrencyCode:           model.CurrencyCode("RSD"),
	}

	p, err := svc.CreatePayment(context.Background(), req)
	require.Nil(t, p)
	require.Error(t, err)
	require.Equal(t, "db error", err.Error())
}

func TestGetAccountPayments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		repo      *fakePaymentRepo
		expectErr bool
		check     func(t *testing.T, payments []model.Payment, total int64)
	}{
		{
			name: "success returns payments",
			repo: &fakePaymentRepo{
				payments: []model.Payment{
					{PaymentID: 1, RecipientName: "Marko Markovic", Transaction: model.Transaction{StartAmount: 5000, StartCurrencyCode: "RSD", Status: model.TransactionCompleted}},
					{PaymentID: 2, RecipientName: "Ana Jovanovic", Transaction: model.Transaction{StartAmount: 1200, StartCurrencyCode: "RSD", Status: model.TransactionProcessing}},
				},
				total: 2,
			},
			check: func(t *testing.T, payments []model.Payment, total int64) {
				require.Len(t, payments, 2)
				require.Equal(t, int64(2), total)
				require.Equal(t, "Marko Markovic", payments[0].RecipientName)
				require.Equal(t, model.TransactionCompleted, payments[0].Transaction.Status)
				require.Equal(t, model.TransactionProcessing, payments[1].Transaction.Status)
			},
		},
		{
			name:      "repo error returns internal error",
			repo:      &fakePaymentRepo{findAccErr: errors.New("db failure")},
			expectErr: true,
		},
		{
			name: "returns empty list when no payments",
			repo: &fakePaymentRepo{payments: []model.Payment{}, total: 0},
			check: func(t *testing.T, payments []model.Payment, total int64) {
				require.Empty(t, payments)
				require.Equal(t, int64(0), total)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := newPaymentService(tt.repo, &fakeTransactionRepo{})
			payments, total, err := svc.GetAccountPayments(context.Background(), "444000112345678911", &dto.PaymentFilters{Page: 1, PageSize: 10})
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.check != nil {
				tt.check(t, payments, total)
			}
		})
	}
}

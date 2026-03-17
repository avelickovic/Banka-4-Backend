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
	createErr   error
	getErr      error
	transaction *model.Transaction
}

func (f *fakeTransactionRepo) Create(ctx context.Context, t *model.Transaction) error {
	if f.createErr != nil {
		return f.createErr
	}
	t.TransactionID = 1
	f.transaction = t
	return nil
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

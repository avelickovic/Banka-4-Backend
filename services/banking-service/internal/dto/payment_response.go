package dto

import (
	"banking-service/internal/model"
	"time"
)

type CreatePaymentResponse struct {
	PaymentID uint `json:"id"`
}

type VerifyPaymentResponse struct {
	PaymentID uint `json:"id"`
}

type PaymentResponse struct {
	PaymentID              uint      `json:"paymentId"`
	RecipientName          string    `json:"recipientName"`
	ReferenceNumber        string    `json:"referenceNumber"`
	PaymentCode            string    `json:"paymentCode"`
	Purpose                string    `json:"purpose"`
	PayerAccountNumber     string    `json:"payerAccountNumber"`
	RecipientAccountNumber string    `json:"recipientAccountNumber"`
	Amount                 float64   `json:"amount"`
	CurrencyCode           string    `json:"currencyCode"`
	Status                 string    `json:"status"`
	CreatedAt              time.Time `json:"createdAt"`
}

func ToPaymentResponse(p *model.Payment) PaymentResponse {
	return PaymentResponse{
		PaymentID:              p.PaymentID,
		RecipientName:          p.RecipientName,
		ReferenceNumber:        p.ReferenceNumber,
		PaymentCode:            p.PaymentCode,
		Purpose:                p.Purpose,
		PayerAccountNumber:     p.Transaction.PayerAccountNumber,
		RecipientAccountNumber: p.Transaction.RecipientAccountNumber,
		Amount:                 p.Transaction.StartAmount,
		CurrencyCode:           string(p.Transaction.StartCurrencyCode),
		Status:                 string(p.Transaction.Status),
		CreatedAt:              p.Transaction.CreatedAt,
	}
}

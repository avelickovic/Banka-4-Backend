package dto

import "banking-service/internal/model"

type PayeeResponse struct {
	PayeeID       uint   `json:"payeeId"`
	ClientID      uint   `json:"clientId"`
	Name          string `json:"name"`
	AccountNumber string `json:"accountNumber"`
}

func ToPayeeResponse(p *model.Payee) PayeeResponse {
	return PayeeResponse{
		PayeeID:       p.PayeeID,
		ClientID:      p.ClientID,
		Name:          p.Name,
		AccountNumber: p.AccountNumber,
	}
}
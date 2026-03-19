package dto

type CreatePayeeRequest struct {
	Name          string `json:"name"          binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
}

type UpdatePayeeRequest struct {
	Name          string `json:"name"          binding:"omitempty"`
	AccountNumber string `json:"accountNumber" binding:"omitempty"`
}
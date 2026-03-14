package dto

type ActivateAccountRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,password"`
}

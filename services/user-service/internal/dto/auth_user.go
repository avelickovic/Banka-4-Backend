package dto

import (
	"common/pkg/auth"
	"common/pkg/permission"
	"user-service/internal/model"
)

type AuthUser struct {
	ID           uint                    `json:"id"`
	IdentityType auth.IdentityType       `json:"identity_type"`
	FirstName    string                  `json:"first_name"`
	LastName     string                  `json:"last_name"`
	Email        string                  `json:"email"`
	Username     string                  `json:"username"`
	Permissions  []permission.Permission `json:"permissions"`
}

func NewAuthUserFromEmployee(identity *model.Identity, employee *model.Employee) *AuthUser {
	return &AuthUser{
		ID:           employee.EmployeeID,
		IdentityType: identity.Type,
		FirstName:    employee.FirstName,
		LastName:     employee.LastName,
		Email:        identity.Email,
		Username:     identity.Username,
		Permissions:  employee.RawPermissions(),
	}
}

func NewAuthUserFromClient(identity *model.Identity, client *model.Client) *AuthUser {
	return &AuthUser{
		ID:           client.ClientID,
		IdentityType: identity.Type,
		FirstName:    client.FirstName,
		LastName:     client.LastName,
		Email:        identity.Email,
		Username:     identity.Username,
		Permissions:  []permission.Permission{},
	}
}

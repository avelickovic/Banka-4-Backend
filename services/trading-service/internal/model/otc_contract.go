package model

import "time"

type OtcContractStatus string

const (
	OtcContractStatusPending  OtcContractStatus = "PENDING"
	OtcContractStatusAccepted OtcContractStatus = "ACCEPTED"
	OtcContractStatusRejected OtcContractStatus = "REJECTED"
)

type OtcContract struct {
	OtcContractID    uint `gorm:"primaryKey;autoIncrement"`
	BuyerID          uint `gorm:"not null;index"`
	SellerID         uint `gorm:"not null;index"`
	AssetID          uint `gorm:"not null;index"`
	Asset            Asset
	Quantity         float64 `gorm:"not null"`
	PricePerUnit     float64 `gorm:"not null"`
	TotalPrice       float64 `gorm:"not null"`
	BankApproved     *bool
	SellerApproved   *bool
	Comment          *string   `gorm:"size:512"`
	ContractNumber   string    `gorm:"not null;uniqueIndex;size:50"`
	FinalizedAt      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
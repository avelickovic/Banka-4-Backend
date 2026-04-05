package model

import (
	"time"
)

type TaxStatus string

const (
	TaxStatusCollected TaxStatus = "COLLECTED"
	TaxStatusFailed    TaxStatus = "FAILED"
)

type AccumulatedTax struct {
	AccumulatedTaxID uint    `gorm:"primaryKey;autoIncrement"`
	AccountNumber    string  `gorm:"not null;uniqueIndex:idx_acc_emp"`
	TaxOwed          float64 `gorm:"not null;default:0"`
	CurrencyCode     string  `gorm:"not null;default:'RSD'"`
	EmployeeID       *uint   `gorm:"uniqueIndex:idx_acc_emp"`
	LastUpdatedAt    time.Time
	LastClearedAt    *time.Time
}

type TaxCollection struct {
	TaxCollectionID   uint    `gorm:"primaryKey;autoIncrement"`
	AccountNumber     string  `gorm:"not null"`
	TaxOwed           float64 `gorm:"not null"`
	CurrencyCode      string  `gorm:"not null;default:'RSD'"`
	EmployeeID        *uint
	Status            TaxStatus `gorm:"type:varchar(20);not null;check:status IN ('COLLECTED','FAILED')"`
	FailureReason     *string   `gorm:"type:text"`
	TaxingPeriodStart time.Time `gorm:"not null"`
	TaxingPeriodEnd   *time.Time
	TriggeredByID     *uint
}

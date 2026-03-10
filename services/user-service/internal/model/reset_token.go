package model

import "time"

type ResetToken struct {
	ID         uint      `gorm:"primaryKey"`
	EmployeeID uint      `gorm:"not null;index"`
	Code       string    `gorm:"not null"`
	ExpiresAt  time.Time `gorm:"not null"`
	//Used       bool      `gorm:"default:false"`
}

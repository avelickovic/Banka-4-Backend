package model

import (
	"time"
)

type Listing struct {
	ListingID   uint      `gorm:"primaryKey;autoIncrement"`
	Ticker      string    `gorm:"not null;uniqueIndex;size:20"`
	Name        string    `gorm:"not null"`
	ExchangeMIC string    `gorm:"not null;size:100"`
	LastRefresh time.Time `gorm:"not null"`
	Price       float64   `gorm:"not null;default:0"`
	Ask         float64   `gorm:"not null;default:0"`
}

type ListingDailyPriceInfo struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	ListingID uint `gorm:"not null;index"`
	Listing   Listing
	Date      time.Time `gorm:"not null;index"`
	Price     float64   `gorm:"not null;default:0"`
	Ask       float64   `gorm:"not null;default:0"`
	Bid       float64   `gorm:"not null;default:0"`
	Change    float64   `gorm:"not null;default:0"`
	Volume    uint      `gorm:"not null;default:0"`
}

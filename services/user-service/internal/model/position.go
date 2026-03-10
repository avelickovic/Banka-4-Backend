package model

type Position struct {
	PositionID uint   `gorm:"primaryKey"`
	Title      string `gorm:"size:100;not null"`
	// Employees  []Employee `gorm:"foreignKey:PositionID"` // 1M veza
}

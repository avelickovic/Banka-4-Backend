package dto

import "github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"

type CreateOrderRequest struct {
	ListingID  uint                 `json:"listingId" binding:"required"`
	AccountID  uint                 `json:"accountId" binding:"required"`
	OrderType  model.OrderType      `json:"orderType" binding:"required,oneof=MARKET LIMIT STOP STOP_LIMIT"`
	Direction  model.OrderDirection `json:"direction" binding:"required,oneof=BUY SELL"`
	Quantity   uint                 `json:"quantity" binding:"required,min=1"`
	LimitValue *float64             `json:"limitValue,omitempty"`
	StopValue  *float64             `json:"stopValue,omitempty"`
	AllOrNone  bool                 `json:"allOrNone"`
	Margin     bool                 `json:"margin"`
}

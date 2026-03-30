package dto

import (
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type OrderResponse struct {
	OrderID           uint                 `json:"orderId"`
	UserID            uint                 `json:"userId"`
	AccountID         uint                 `json:"accountId"`
	ListingID         uint                 `json:"listingId"`
	Ticker            string               `json:"ticker"`
	ListingName       string               `json:"listingName"`
	OrderType         model.OrderType      `json:"orderType"`
	Direction         model.OrderDirection `json:"direction"`
	Quantity          uint                 `json:"quantity"`
	ContractSize      float64              `json:"contractSize"`
	PricePerUnit      *float64             `json:"pricePerUnit"`
	LimitValue        *float64             `json:"limitValue,omitempty"`
	StopValue         *float64             `json:"stopValue,omitempty"`
	AllOrNone         bool                 `json:"allOrNone"`
	Margin            bool                 `json:"margin"`
	Status            model.OrderStatus    `json:"status"`
	ApprovedBy        *uint                `json:"approvedBy,omitempty"`
	IsDone            bool                 `json:"isDone"`
	AfterHours        bool                 `json:"afterHours"`
	RemainingPortions uint                 `json:"remainingPortions"`
	CreatedAt         time.Time            `json:"createdAt"`
	UpdatedAt         time.Time            `json:"updatedAt"`
}

func ToOrderResponse(o model.Order) OrderResponse {
	return OrderResponse{
		OrderID:           o.OrderID,
		UserID:            o.UserID,
		AccountID:         o.AccountID,
		ListingID:         o.ListingID,
		Ticker:            o.Listing.Ticker,
		ListingName:       o.Listing.Name,
		OrderType:         o.OrderType,
		Direction:         o.Direction,
		Quantity:          o.Quantity,
		ContractSize:      o.ContractSize,
		PricePerUnit:      o.PricePerUnit,
		LimitValue:        o.LimitValue,
		StopValue:         o.StopValue,
		AllOrNone:         o.AllOrNone,
		Margin:            o.Margin,
		Status:            o.Status,
		ApprovedBy:        o.ApprovedBy,
		IsDone:            o.IsDone,
		AfterHours:        o.AfterHours,
		RemainingPortions: o.RemainingPortions(),
		CreatedAt:         o.CreatedAt,
		UpdatedAt:         o.UpdatedAt,
	}
}

func ToOrderResponseList(orders []model.Order) []OrderResponse {
	result := make([]OrderResponse, len(orders))
	for i, o := range orders {
		result[i] = ToOrderResponse(o)
	}
	return result
}

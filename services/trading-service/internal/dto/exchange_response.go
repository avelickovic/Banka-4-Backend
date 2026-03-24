package dto

import "github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"

type ExchangeResponse struct {
	ExchangeID     uint   `json:"exchangeId"`
	Name           string `json:"name"`
	Acronym        string `json:"acronym"`
	MicCode        string `json:"micCode"`
	Polity         string `json:"polity"`
	Currency       string `json:"currency"`
	TimeZone       int    `json:"timeZone"`
	OpenTime       string `json:"openTime"`
	CloseTime      string `json:"closeTime"`
	TradingEnabled bool   `json:"tradingEnabled"`
}

func ToExchangeResponse(e model.Exchange) ExchangeResponse {
	return ExchangeResponse{
		ExchangeID:     e.ExchangeID,
		Name:           e.Name,
		Acronym:        e.Acronym,
		MicCode:        e.MicCode,
		Polity:         e.Polity,
		Currency:       e.Currency,
		TimeZone:       e.TimeZone,
		OpenTime:       e.OpenTime,
		CloseTime:      e.CloseTime,
		TradingEnabled: e.TradingEnabled,
	}
}

func ToExchangeResponseList(exchanges []model.Exchange) []ExchangeResponse {
	result := make([]ExchangeResponse, len(exchanges))
	for i, e := range exchanges {
		result[i] = ToExchangeResponse(e)
	}
	return result
}

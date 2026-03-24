package dto

type ListExchangesQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

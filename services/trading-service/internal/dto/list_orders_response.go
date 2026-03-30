package dto

type ListOrdersResponse struct {
	Data     []OrderResponse `json:"data"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

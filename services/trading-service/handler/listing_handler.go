package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/service"
)

type ListingHandler struct {
	svc *service.ListingService
}

func NewListingHandler(svc *service.ListingService) *ListingHandler {
	return &ListingHandler{svc: svc}
}

// GetStocks godoc
// @Summary List stocks
// @Tags listings
// @Produce json
// @Param search query string false "Search by ticker or name"
// @Param exchange query string false "Exchange MIC prefix"
// @Param price_min query number false "Min price"
// @Param price_max query number false "Max price"
// @Param ask_min query number false "Min ask"
// @Param ask_max query number false "Max ask"
// @Param bid_min query number false "Min bid"
// @Param bid_max query number false "Max bid"
// @Param volume_min query integer false "Min volume"
// @Param volume_max query integer false "Max volume"
// @Param sort_by query string false "price|volume|maintenance_margin"
// @Param sort_dir query string false "asc|desc"
// @Param page query integer false "Page"
// @Param page_size query integer false "Page size"
// @Success 200 {object} dto.PaginatedResponse[dto.StockResponse]
// @Router /api/listings/stocks [get]
func (h *ListingHandler) GetStocks(c *gin.Context) {
	var q dto.ListingQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}
	q.Normalize()

	result, err := h.svc.GetStocks(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetFutures godoc
// @Summary List futures contracts
// @Tags listings
// @Produce json
// @Param search query string false "Search by ticker or name"
// @Param exchange query string false "Exchange MIC prefix"
// @Param price_min query number false "Min price"
// @Param price_max query number false "Max price"
// @Param ask_min query number false "Min ask"
// @Param ask_max query number false "Max ask"
// @Param bid_min query number false "Min bid"
// @Param bid_max query number false "Max bid"
// @Param volume_min query integer false "Min volume"
// @Param volume_max query integer false "Max volume"
// @Param settlement_date query string false "Filter by settlement date (YYYY-MM-DD)"
// @Param sort_by query string false "price|volume|maintenance_margin"
// @Param sort_dir query string false "asc|desc"
// @Param page query integer false "Page"
// @Param page_size query integer false "Page size"
// @Success 200 {object} dto.PaginatedResponse[dto.FuturesResponse]
// @Router /api/listings/futures [get]
func (h *ListingHandler) GetFutures(c *gin.Context) {
	var q dto.ListingQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}
	q.Normalize()

	result, err := h.svc.GetFutures(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetForex godoc
// @Summary List forex pairs
// @Tags listings
// @Produce json
// @Param search query string false "Search by base/quote"
// @Param price_min query number false "Min rate"
// @Param price_max query number false "Max rate"
// @Param sort_dir query string false "asc|desc"
// @Param page query integer false "Page"
// @Param page_size query integer false "Page size"
// @Success 200 {object} dto.PaginatedResponse[dto.ForexResponse]
// @Router /api/listings/forex [get]
func (h *ListingHandler) GetForex(c *gin.Context) {
	var q dto.ListingQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}
	q.Normalize()

	result, err := h.svc.GetForex(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetOptions godoc
// @Summary List options
// @Tags listings
// @Produce json
// @Param search query string false "Search by ticker or name"
// @Param exchange query string false "Exchange MIC prefix"
// @Param price_min query number false "Min price"
// @Param price_max query number false "Max price"
// @Param ask_min query number false "Min ask"
// @Param ask_max query number false "Max ask"
// @Param bid_min query number false "Min bid"
// @Param bid_max query number false "Max bid"
// @Param volume_min query integer false "Min volume"
// @Param volume_max query integer false "Max volume"
// @Param settlement_date query string false "Filter by settlement date (YYYY-MM-DD)"
// @Param sort_by query string false "price|volume|maintenance_margin"
// @Param sort_dir query string false "asc|desc"
// @Param page query integer false "Page"
// @Param page_size query integer false "Page size"
// @Success 200 {object} dto.PaginatedResponse[dto.OptionResponse]
// @Router /api/listings/options [get]
func (h *ListingHandler) GetOptions(c *gin.Context) {
	var q dto.ListingQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}
	q.Normalize()

	result, err := h.svc.GetOptions(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, result)
}

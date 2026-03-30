package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Creates a buy or sell order for a listing
// @Tags orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderRequest true "Order details"
// @Success 201 {object} dto.OrderResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /api/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToOrderResponse(*order))
}

// ApproveOrder godoc
// @Summary Approve a pending order
// @Description Supervisor approves a pending order
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /api/orders/{id}/approve [patch]
func (h *OrderHandler) ApproveOrder(c *gin.Context) {
	orderID, err := parseOrderID(c)
	if err != nil {
		c.Error(err)
		return
	}

	order, err := h.service.ApproveOrder(c.Request.Context(), orderID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderResponse(*order))
}

// DeclineOrder godoc
// @Summary Decline a pending order
// @Description Supervisor declines a pending order
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /api/orders/{id}/decline [patch]
func (h *OrderHandler) DeclineOrder(c *gin.Context) {
	orderID, err := parseOrderID(c)
	if err != nil {
		c.Error(err)
		return
	}

	order, err := h.service.DeclineOrder(c.Request.Context(), orderID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderResponse(*order))
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel a pending or approved order that hasn't been fully executed
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /api/orders/{id}/cancel [patch]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID, err := parseOrderID(c)
	if err != nil {
		c.Error(err)
		return
	}

	order, err := h.service.CancelOrder(c.Request.Context(), orderID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderResponse(*order))
}

func parseOrderID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, errors.BadRequestErr("invalid order id")
	}

	return uint(id), nil
}

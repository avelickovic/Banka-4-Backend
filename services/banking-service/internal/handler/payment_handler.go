package handler

import (
	"banking-service/internal/dto"
	"banking-service/internal/model"
	"banking-service/internal/repository"
	"banking-service/internal/service"
	"common/pkg/errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) GetPayments(c *gin.Context) {
	filter := repository.PaymentFilter{}

	if v := c.Query("date_from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			c.Error(errors.BadRequestErr("invalid date_from format, use RFC3339"))
			return
		}
		filter.DateFrom = &t
	}

	if v := c.Query("date_to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			c.Error(errors.BadRequestErr("invalid date_to format, use RFC3339"))
			return
		}
		filter.DateTo = &t
	}

	if v := c.Query("amount_min"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			c.Error(errors.BadRequestErr("invalid amount_min"))
			return
		}
		filter.AmountMin = &f
	}

	if v := c.Query("amount_max"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			c.Error(errors.BadRequestErr("invalid amount_max"))
			return
		}
		filter.AmountMax = &f
	}

	if v := c.Query("status"); v != "" {
		s := model.TransactionStatus(v)
		filter.Status = &s
	}

	payments, err := h.service.GetFilteredPayments(c.Request.Context(), filter)
	if err != nil {
		c.Error(err)
		return
	}

	response := make([]dto.PaymentResponse, len(payments))
	for i, p := range payments {
		response[i] = dto.ToPaymentResponse(&p)
	}

	c.JSON(http.StatusOK, response)
}


func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}


	payment, err := h.service.CreatePayment(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.CreatePaymentResponse{
		PaymentID: payment.PaymentID,
	})
}

func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequestErr("invalid payment id"))
		return
	}

	var req dto.VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}

	payment, err := h.service.VerifyPayment(c.Request.Context(), uint(id), req.Code)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.VerifyPaymentResponse{
		PaymentID: payment.PaymentID,
	})
}

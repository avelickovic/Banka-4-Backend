package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/service"
)

type CompanyHandler struct {
	service *service.CompanyService
}

func NewCompanyHandler(service *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

// GetWorkCodes godoc
// @Summary List all work codes
// @Description Returns all available company work codes for frontend selection.
// @Tags companies
// @Produce json
// @Success 200 {array} dto.WorkCodeResponse
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Security BearerAuth
// @Router /api/companies/work-codes [get]
func (h *CompanyHandler) GetWorkCodes(c *gin.Context) {
	workCodes, err := h.service.GetWorkCodes(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToWorkCodeResponses(workCodes))
}

func (h *CompanyHandler) Create(c *gin.Context) {
	var req dto.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}

	company, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToCompanyResponse(company))
}

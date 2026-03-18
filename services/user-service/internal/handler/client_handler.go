package handler

import (
	"net/http"

	"common/pkg/errors"
	"user-service/internal/dto"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	service *service.ClientService
}

func NewClientHandler(service *service.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

// Register godoc
// @Summary Register a new client
// @Description Creates a new client account and sends an activation email
// @Tags clients
// @Accept json
// @Produce json
// @Param client body dto.CreateClientRequest true "Client registration data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} errors.AppError
// @Failure 409 {object} errors.AppError
// @Failure 503 {object} errors.AppError
// @Router /api/clients/register [post]
func (h *ClientHandler) Register(c *gin.Context) {
	var req dto.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}

	_, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful. Please check your email to activate your account."})
}

// GetMobileSecret godoc
// @Summary Get mobile verification secret
// @Description Returns TOTP secret for currently authenticated client
// @Tags clients
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.MobileSecretResponse
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Router /api/secret-mobile [get]
func (h *ClientHandler) GetMobileSecret(c *gin.Context) {
	secret, err := h.service.GetMobileVerificationSecret(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.MobileSecretResponse{Secret: secret})
}

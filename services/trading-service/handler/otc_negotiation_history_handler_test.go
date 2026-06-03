package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	app_errors "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type mockOtcHistoryService struct {
	mock.Mock
}

func (m *mockOtcHistoryService) CreateNegotiationHistory(ctx context.Context, offerID uint, oldOffer, newOffer *model.OtcOffer, modifiedBy uint) error {
	args := m.Called(ctx, offerID, oldOffer, newOffer, modifiedBy)
	return args.Error(0)
}

func (m *mockOtcHistoryService) GetNegotiationHistory(ctx context.Context, offerID uint, statusFilter string, dateFrom, dateTo *time.Time, counterpartyFilter uint) ([]*model.OtcNegotiationHistory, error) {
	args := m.Called(ctx, offerID, statusFilter, dateFrom, dateTo, counterpartyFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.OtcNegotiationHistory), args.Error(1)
}

// --- Test Setup ---

func setupHandlerTest() (*gin.Engine, *mockOtcHistoryService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(mockOtcHistoryService)
	historyHandler := NewOtcNegotiationHistoryHandler(mockService)
	router.GET("/api/otc/offers/:id/history", historyHandler.GetNegotiationHistory)
	return router, mockService
}

// --- Tests ---

func TestGetNegotiationHistoryHandler_HappyPath(t *testing.T) {
	router, mockService := setupHandlerTest()

	expectedHistory := []*model.OtcNegotiationHistory{
		{OtcNegotiationHistoryID: 1, OtcOfferID: 1},
	}
	mockService.On("GetNegotiationHistory", mock.Anything, uint(1), "ACCEPTED", mock.Anything, mock.Anything, uint(123)).Return(expectedHistory, nil)

	req, _ := http.NewRequest("GET", "/api/otc/offers/1/history?status=ACCEPTED&counterparty=123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody []*model.OtcNegotiationHistory
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedHistory, responseBody)
	mockService.AssertExpectations(t)
}

func TestGetNegotiationHistoryHandler_ServiceError(t *testing.T) {
	router, mockService := setupHandlerTest()

	expectedError := app_errors.NotFoundErr("offer not found")
	mockService.On("GetNegotiationHistory", mock.Anything, uint(1), "", mock.Anything, mock.Anything, uint(0)).Return(nil, expectedError)

	// We need to add the error handler middleware to the router for this test
	router.Use(app_errors.ErrorHandler())
	req, _ := http.NewRequest("GET", "/api/otc/offers/1/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse app_errors.AppError
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Not Found", errorResponse.Type)
	assert.Equal(t, "offer not found", errorResponse.Message)
	mockService.AssertExpectations(t)
}

func TestGetNegotiationHistoryHandler_InvalidUrlParams(t *testing.T) {
	tt := []struct {
		name       string
		url        string
		expectedMsg string
	}{
		{"Invalid offer ID", "/api/otc/offers/abc/history", "invalid offer id"},
		{"Invalid counterparty", "/api/otc/offers/1/history?counterparty=xyz", "invalid counterparty"},
		{"Invalid from date", "/api/otc/offers/1/history?from=not-a-date", "invalid from date"},
		{"Invalid to date", "/api/otc/offers/1/history?to=not-a-date", "invalid to date"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router, _ := setupHandlerTest()
			router.Use(app_errors.ErrorHandler())

			req, _ := http.NewRequest("GET", tc.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			var errorResponse app_errors.AppError
			err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
			assert.NoError(t, err)
			assert.Equal(t, "Bad Request", errorResponse.Type)
			assert.Contains(t, errorResponse.Message, tc.expectedMsg)
		})
	}
}

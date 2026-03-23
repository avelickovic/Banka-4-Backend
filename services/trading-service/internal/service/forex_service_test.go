package service

import (
	"context"
	"testing"
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/client"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// --- Mock client za testove ---
type mockExchangeClient struct {
	data *client.ExchangeRateAPIResponse
}

func (m *mockExchangeClient) FetchRates(ctx context.Context) (*client.ExchangeRateAPIResponse, error) {
	return m.data, nil
}

// --- Helper funkcija za in-memory CGO-free DB (unikatna baza po testu) ---
func setupTestDB(t *testing.T) *gorm.DB {
	// Unikatno ime baze po testu → izbegava deljenje između testova
	dsn := "file:testdb_" + time.Now().Format("150405.000") + "?mode=memory&_pragma=foreign_keys(1)"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&model.ForexPair{}); err != nil {
		t.Fatal(err)
	}

	return db
}

// --- Test za refreshFromAPI ---
func TestRefreshFromAPI(t *testing.T) {
	db := setupTestDB(t)

	mockResp := &client.ExchangeRateAPIResponse{
		BaseCode:           "RSD",
		TimeLastUpdateUnix: time.Now().Unix(),
		TimeNextUpdateUnix: time.Now().Add(time.Hour).Unix(),
		ConversionRates: map[string]float64{
			"USD": 0.0085,
			"EUR": 0.0080,
		},
	}

	mockClient := &mockExchangeClient{data: mockResp}
	service := NewForexService(db, mockClient)

	if err := service.refreshFromAPI(context.Background()); err != nil {
		t.Fatalf("refreshFromAPI failed: %v", err)
	}

	var pairs []model.ForexPair
	if err := db.Find(&pairs).Error; err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if len(pairs) != 2 {
		t.Fatalf("expected 2 forex pairs, got %d", len(pairs))
	}

	for _, pair := range pairs {
		if pair.Base != "RSD" {
			t.Errorf("expected base RSD, got %s", pair.Base)
		}
		if pair.Quote != "USD" && pair.Quote != "EUR" {
			t.Errorf("unexpected quote %s", pair.Quote)
		}
	}
}

// --- Test za Initialize i seeding ---
func TestInitialize_SeedsDB(t *testing.T) {
	db := setupTestDB(t)

	mockResp := &client.ExchangeRateAPIResponse{
		BaseCode:           "RSD",
		TimeLastUpdateUnix: time.Now().Unix(),
		TimeNextUpdateUnix: time.Now().Add(time.Hour).Unix(),
		ConversionRates: map[string]float64{
			"USD": 0.0085,
		},
	}

	mockClient := &mockExchangeClient{data: mockResp}
	service := NewForexService(db, mockClient)

	// DB prazna → Initialize seeduje
	service.Initialize(context.Background())

	var count int64
	if err := db.Model(&model.ForexPair{}).Count(&count).Error; err != nil {
		t.Fatalf("count query failed: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 forex pair, got %d", count)
	}

	// Sada DB ima 1 par, dodajemo novu valutu u mock client
	mockClient.data.ConversionRates["EUR"] = 0.0080

	// Initialize se poziva, ali DB već nije prazna → ne ubacuje EUR
	service.Initialize(context.Background())

	if err := db.Model(&model.ForexPair{}).Count(&count).Error; err != nil {
		t.Fatalf("count query failed: %v", err)
	}

	if count != 1 { // Očekujemo da i dalje bude 1
		t.Fatalf("expected count still 1, got %d", count)
	}
}

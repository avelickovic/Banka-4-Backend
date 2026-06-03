package service

import (
	"github.com/vesic/banka-4-backend/services/trading-service/internal/model"
	"github.com/vesic/banka-4-backend/services/trading-service/internal/repository"
	"time"
)

type OtcNegotiationHistoryService interface {
	CreateNegotiationHistory(
		offerID uint,
		oldOffer, newOffer *model.OtcOffer,
		modifiedBy uint,
	) error
	GetNegotiationHistory(
		offerID uint,
		statusFilter string,
		dateFrom, dateTo *time.Time,
		counterpartyFilter uint,
	) ([]*model.OtcNegotiationHistory, error)
}

type otcNegotiationHistoryService struct {
	otcOfferRepo          repository.OtcOfferRepository
	otcNegotiationHistoryRepo repository.OtcNegotiationHistoryRepository
}

func NewOtcNegotiationHistoryService(
	otcOfferRepo repository.OtcOfferRepository,
	otcNegotiationHistoryRepo repository.OtcNegotiationHistoryRepository,
) OtcNegotiationHistoryService {
	return &otcNegotiationHistoryService{
		otcOfferRepo:          otcOfferRepo,
		otcNegotiationHistoryRepo: otcNegotiationHistoryRepo,
	}
}

func (s *otcNegotiationHistoryService) CreateNegotiationHistory(
	offerID uint,
	oldOffer, newOffer *model.OtcOffer,
	modifiedBy uint,
) error {
	history := &model.OtcNegotiationHistory{
		OtcOfferID:          offerID,
		OldAmount:           oldOffer.Amount,
		NewAmount:           newOffer.Amount,
		OldPricePerStockRSD: oldOffer.PricePerStockRSD,
		NewPricePerStockRSD: newOffer.PricePerStockRSD,
		OldPremiumRSD:       oldOffer.PremiumRSD,
		NewPremiumRSD:       newOffer.PremiumRSD,
		OldSettlementDate:   oldOffer.SettlementDate,
		NewSettlementDate:   newOffer.SettlementDate,
		Timestamp:           time.Now(),
		ModifiedBy:          modifiedBy,
	}
	return s.otcNegotiationHistoryRepo.Create(history)
}

func (s *otcNegotiationHistoryService) GetNegotiationHistory(
	offerID uint,
	statusFilter string,
	dateFrom, dateTo *time.Time,
	counterpartyFilter uint,
) ([]*model.OtcNegotiationHistory, error) {
	offer, err := s.otcOfferRepo.FindByID(offerID)
	if err != nil {
		return nil, err
	}

	if offer.Status == model.OtcOfferStatusActive {
		return nil, nil // Or an error indicating the negotiation is still active
	}

	// Apply filters
	if statusFilter != "" && string(offer.Status) != statusFilter {
		return nil, nil
	}
	if dateFrom != nil && offer.LastModified.Before(*dateFrom) {
		return nil, nil
	}
	if dateTo != nil && offer.LastModified.After(*dateTo) {
		return nil, nil
	}
	if counterpartyFilter != 0 && (offer.BuyerID != counterpartyFilter && offer.SellerID != counterpartyFilter) {
		return nil, nil
	}

	return s.otcNegotiationHistoryRepo.GetByOfferID(offerID)
}

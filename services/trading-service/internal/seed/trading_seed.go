package seed

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

const (
	finnhubBaseURL = "https://finnhub.io/api/v1"
	rateLimit      = time.Second
)

// --- Finnhub response structs ---

type FinnhubSymbol struct {
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	MIC         string `json:"mic"`
}

type FinnhubProfile struct {
	Name             string  `json:"name"`
	Exchange         string  `json:"exchange"`
	ShareOutstanding float64 `json:"shareOutstanding"`
	Ticker           string  `json:"ticker"`
}

type FinnhubQuote struct {
	CurrentPrice float64 `json:"c"`
	High         float64 `json:"h"`
}

type FinnhubBasicFinancials struct {
	Metric struct {
		DividendYieldIndicatedAnnual float64 `json:"dividendYieldIndicatedAnnual"`
	} `json:"metric"`
}

// --- Main seed function ---
const maxCallsPerMinute = 55 //60 is max per minute

func SeedStocks(db *gorm.DB, apiKey string, limit int) error {
	get := func(path string, out interface{}) error {
		separator := "?"
		for _, c := range path {
			if c == '?' {
				separator = "&"
				break
			}
		}
		url := fmt.Sprintf("%s%s%stoken=%s", finnhubBaseURL, path, separator, apiKey)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("finnhub returned status %d", resp.StatusCode)
		}
		return json.NewDecoder(resp.Body).Decode(out)
	}

	// 1. Fetch list of US stocks (1 call)
	var symbols []FinnhubSymbol
	if err := get("/stock/symbol?exchange=US", &symbols); err != nil {
		return fmt.Errorf("failed to fetch symbol list: %w", err)
	}

	log.Printf("Fetched %d symbols, seeding up to %d", len(symbols), limit)

	callsThisMinute := 1 // already used 1 for symbol list
	minuteStart := time.Now()

	count := 0
	for _, sym := range symbols {
		if count >= limit {
			break
		}

		ticker := sym.Symbol
		if containsDot(ticker) {
			continue
		}

		// Check if we need to wait for the next minute window
		// Each stock uses 3 calls, so check before starting a new stock
		if callsThisMinute+3 > maxCallsPerMinute {
			elapsed := time.Since(minuteStart)
			if elapsed < time.Minute {
				wait := time.Minute - elapsed
				log.Printf("Rate limit window reached, waiting %s...", wait.Round(time.Second))
				time.Sleep(wait)
			}
			callsThisMinute = 0
			minuteStart = time.Now()
		}

		// 2. profile
		var profile FinnhubProfile
		if err := get(fmt.Sprintf("/stock/profile2?symbol=%s", ticker), &profile); err != nil {
			log.Printf("skipping %s: profile error: %v", ticker, err)
			callsThisMinute++
			continue
		}
		callsThisMinute++
		if profile.Name == "" {
			log.Printf("skipping %s: empty profile", ticker)
			continue
		}

		// 3. quote
		var quote FinnhubQuote
		if err := get(fmt.Sprintf("/quote?symbol=%s", ticker), &quote); err != nil {
			log.Printf("skipping %s: quote error: %v", ticker, err)
			callsThisMinute++
			continue
		}
		callsThisMinute++
		if quote.CurrentPrice == 0 {
			log.Printf("skipping %s: no price data", ticker)
			continue
		}

		// 4. financials
		var financials FinnhubBasicFinancials
		if err := get(fmt.Sprintf("/stock/metric?symbol=%s&metric=all", ticker), &financials); err != nil {
			log.Printf("skipping %s: financials error: %v", ticker, err)
			callsThisMinute++
			continue
		}
		callsThisMinute++

		// 5. Upsert Listing
		listing := model.Listing{
			Ticker:      ticker,
			Name:        profile.Name,
			ExchangeMIC: profile.Exchange,
			LastRefresh: time.Now(),
			Price:       quote.CurrentPrice,
			Ask:         quote.High,
		}
		if err := db.Where(model.Listing{Ticker: ticker}).
			Assign(listing).
			FirstOrCreate(&listing).Error; err != nil {
			log.Printf("skipping %s: db listing error: %v", ticker, err)
			continue
		}

		// 6. Upsert Stock
		stock := model.Stock{
			ListingID:         listing.ListingID,
			OutstandingShares: profile.ShareOutstanding,
			DividendYield:     financials.Metric.DividendYieldIndicatedAnnual,
		}
		if err := db.Where(model.Stock{ListingID: listing.ListingID}).
			Assign(stock).
			FirstOrCreate(&stock).Error; err != nil {
			log.Printf("skipping %s: db stock error: %v", ticker, err)
			continue
		}

		count++
		log.Printf("[%d/%d] seeded %s", count, limit, ticker)
	}

	log.Printf("Done. Seeded %d stocks.", count)
	return nil
}

func containsDot(s string) bool {
	for _, c := range s {
		if c == '.' {
			return true
		}
	}
	return false
}

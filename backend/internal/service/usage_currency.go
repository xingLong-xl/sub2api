package service

import (
	"context"
	"math"
	"strconv"
	"strings"
)

const (
	UsageCurrencyCNY  = "CNY"
	UsageCurrencyUSD  = "USD"
	DefaultUSDCNYRate = 7.2
)

func normalizeUsageCurrency(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case UsageCurrencyUSD:
		return UsageCurrencyUSD
	default:
		return UsageCurrencyCNY
	}
}

func normalizeUSDCNYRate(value float64) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) || value <= 0 {
		return DefaultUSDCNYRate
	}
	return value
}

func parsePositiveFloatDefault(raw string, fallback float64) float64 {
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return fallback
	}
	return normalizeUSDCNYRate(value)
}

// GetUsageCurrency returns the display currency for usage/quota amounts and the USD/CNY rate.
// Stored amounts are canonical CNY; USD is display-only and must not change billing values.
func (s *SettingService) GetUsageCurrency(ctx context.Context) (string, float64) {
	currency := UsageCurrencyCNY
	rate := DefaultUSDCNYRate
	if s == nil || s.settingRepo == nil {
		return currency, rate
	}
	if value, err := s.settingRepo.GetValue(ctx, SettingKeyUsageCurrency); err == nil {
		currency = normalizeUsageCurrency(value)
	}
	if value, err := s.settingRepo.GetValue(ctx, SettingKeyUSDCNYRate); err == nil {
		rate = parsePositiveFloatDefault(value, DefaultUSDCNYRate)
	}
	return currency, rate
}

// ConvertCNYToUsageDisplay converts canonical CNY into the configured display currency.
func ConvertCNYToUsageDisplay(amountCNY float64, currency string, usdCNYRate float64) float64 {
	if normalizeUsageCurrency(currency) != UsageCurrencyUSD {
		return amountCNY
	}
	rate := normalizeUSDCNYRate(usdCNYRate)
	return amountCNY / rate
}

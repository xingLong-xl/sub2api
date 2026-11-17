package service

import "testing"

func TestNormalizeUsageCurrency(t *testing.T) {
	if got := normalizeUsageCurrency("usd"); got != UsageCurrencyUSD {
		t.Fatalf("normalizeUsageCurrency(usd) = %q", got)
	}
	if got := normalizeUsageCurrency(""); got != UsageCurrencyCNY {
		t.Fatalf("normalizeUsageCurrency(empty) = %q", got)
	}
}

func TestConvertCNYToUsageDisplay(t *testing.T) {
	if got := ConvertCNYToUsageDisplay(72, UsageCurrencyCNY, 7.2); got != 72 {
		t.Fatalf("CNY conversion = %v, want 72", got)
	}
	if got := ConvertCNYToUsageDisplay(72, UsageCurrencyUSD, 7.2); got != 10 {
		t.Fatalf("USD conversion = %v, want 10", got)
	}
	if got := ConvertCNYToUsageDisplay(72, UsageCurrencyUSD, 0); got != 10 {
		t.Fatalf("USD conversion with invalid rate = %v, want 10", got)
	}
}

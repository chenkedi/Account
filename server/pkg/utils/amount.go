package utils

import "math"

func RoundToCents(amount float64) float64 {
	return math.Round(amount*100) / 100
}

func IsValidCurrency(code string) bool {
	validCurrencies := map[string]bool{
		"CNY": true,
		"USD": true,
		"EUR": true,
		"GBP": true,
		"JPY": true,
		"HKD": true,
	}
	return validCurrencies[code]
}

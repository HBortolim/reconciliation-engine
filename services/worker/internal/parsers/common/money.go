package common

import "fmt"

// CentsFromReais converts a float64 amount in reais to an integer amount in cents.
func CentsFromReais(reais float64) int64 {
	return int64(reais * 100)
}

// ReaisFromCents converts an integer amount in cents to a float64 amount in reais.
// This is for presentation purposes only.
func ReaisFromCents(cents int64) float64 {
	return float64(cents) / 100.0
}

// AddCents adds two amounts in cents and returns the result.
func AddCents(a, b int64) int64 {
	return a + b
}

// SubtractCents subtracts amount b from amount a (both in cents) and returns the result.
func SubtractCents(a, b int64) int64 {
	return a - b
}

// FormatCents formats an amount in cents as a currency string (R$ X,XX).
func FormatCents(cents int64) string {
	return fmt.Sprintf("R$ %.2f", ReaisFromCents(cents))
}

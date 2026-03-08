package common

import "fmt"

// CentavosFromReais converts a float64 amount in reais to an integer amount in centavos.
func CentavosFromReais(reais float64) int64 {
	return int64(reais * 100)
}

// ReaisFromCentavos converts an integer amount in centavos to a float64 amount in reais.
// This is for presentation purposes only.
func ReaisFromCentavos(centavos int64) float64 {
	return float64(centavos) / 100.0
}

// AddCentavos adds two amounts in centavos and returns the result.
func AddCentavos(a, b int64) int64 {
	return a + b
}

// SubtractCentavos subtracts amount b from amount a (both in centavos) and returns the result.
func SubtractCentavos(a, b int64) int64 {
	return a - b
}

// FormatCentavos formats an amount in centavos as a currency string (R$ X,XX).
func FormatCentavos(centavos int64) string {
	return fmt.Sprintf("R$ %.2f", ReaisFromCentavos(centavos))
}

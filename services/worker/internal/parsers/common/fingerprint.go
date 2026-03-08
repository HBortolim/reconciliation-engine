package common

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// ComputeFingerprint generates a SHA-256 hash fingerprint based on key transaction fields.
// The fingerprint is computed from: source_type + amount + transaction_date + external_id + counterparty_document.
func ComputeFingerprint(sourceType SourceType, amountCentavos int64, transactionDate time.Time, externalID, counterpartyDocument string) string {
	data := fmt.Sprintf("%s|%d|%s|%s|%s",
		sourceType,
		amountCentavos,
		transactionDate.Format(time.RFC3339),
		externalID,
		counterpartyDocument,
	)

	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// ValidateFingerprint checks if a given fingerprint matches the computed one.
func ValidateFingerprint(record *TransactionRecord) bool {
	computed := ComputeFingerprint(
		record.SourceType,
		record.AmountCentavos,
		record.TransactionDate,
		record.ExternalID,
		record.CounterpartyDocument,
	)
	return computed == record.FingerprintHash
}

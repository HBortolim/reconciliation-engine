package common

import (
	"encoding/json"
	"io"
	"time"
)

// SourceType represents the type of transaction source.
type SourceType string

const (
	SourceTypePIX              SourceType = "PIX"
	SourceTypeBOLETO           SourceType = "BOLETO"
	SourceTypeCardCredit       SourceType = "CARD_CREDIT"
	SourceTypeCardDebit        SourceType = "CARD_DEBIT"
	SourceTypeTED              SourceType = "TED"
	SourceTypeDOC              SourceType = "DOC"
	SourceTypeDebitoAutomatico SourceType = "DEBITO_AUTOMATICO"
)

// TransactionRecord represents a unified transaction record across all formats.
type TransactionRecord struct {
	// Core identification
	ID                      string           `json:"id"`
	SourceType              SourceType       `json:"source_type"`
	ExternalID              string           `json:"external_id"`
	NSU                     string           `json:"nsu,omitempty"`
	E2EID                   string           `json:"e2e_id,omitempty"`
	NossoNumero             string           `json:"nosso_numero,omitempty"`
	CounterpartyDocument    string           `json:"counterparty_document"`
	CounterpartyName        string           `json:"counterparty_name,omitempty"`

	// Financial amounts (in centavos)
	AmountCentavos       int64   `json:"amount_centavos"`
	FeeCentavos          int64   `json:"fee_centavos,omitempty"`
	NetAmountCentavos    int64   `json:"net_amount_centavos"`
	DiscountCentavos     int64   `json:"discount_centavos,omitempty"`

	// Dates
	ExpectedSettlementDate time.Time  `json:"expected_settlement_date"`
	ActualSettlementDate   *time.Time `json:"actual_settlement_date,omitempty"`
	TransactionDate        time.Time  `json:"transaction_date"`

	// Fingerprint and raw data
	FingerprintHash string          `json:"fingerprint_hash"`
	RawData         json.RawMessage `json:"raw_data"`
	SourceFile      string          `json:"source_file"`
	ParsedAt        time.Time       `json:"parsed_at"`
}

// Parser is the interface that all parsers must implement.
type Parser interface {
	Parse(reader io.Reader, filename string) ([]TransactionRecord, error)
}

package stone

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// Parser implements the common.Parser interface for Stone JSON API format.
// Stone provides settlement data through JSON API responses or exported JSON files.
type Parser struct{}

// NewParser creates a new Stone parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads a Stone JSON settlement file and returns a slice of TransactionRecord.
// Expected JSON structure: { "transactions": [ { "id", "amount", "status", "created_at", "customer", "nsu" } ] }
func (p *Parser) Parse(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	var payload struct {
		Transactions []StoneTransaction `json:"transactions"`
	}
	if err := json.NewDecoder(reader).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	records := make([]common.TransactionRecord, 0, len(payload.Transactions))
	for i := range payload.Transactions {
		record := p.parseStoneTransaction(&payload.Transactions[i], filename)
		records = append(records, *record)
	}

	return records, nil
}

// StoneTransaction represents a Stone transaction in JSON format.
type StoneTransaction struct {
	ID        string `json:"id"`
	Amount    int64  `json:"amount"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	Customer  struct {
		Document string `json:"document"`
	} `json:"customer"`
	NSU string `json:"nsu,omitempty"`
}

// parseStoneTransaction converts a Stone transaction to a TransactionRecord.
func (p *Parser) parseStoneTransaction(txn *StoneTransaction, filename string) *common.TransactionRecord {
	transDate := parseDate(txn.CreatedAt)
	if transDate.IsZero() {
		transDate = time.Now()
	}

	return &common.TransactionRecord{
		ID:                     txn.ID,
		SourceType:             common.SourceTypeCardCredit,
		NSU:                    txn.NSU,
		ExternalID:             txn.ID,
		CounterpartyDocument:   txn.Customer.Document,
		AmountCentavos:         txn.Amount,
		NetAmountCentavos:      txn.Amount,
		TransactionDate:        transDate,
		ExpectedSettlementDate: transDate.AddDate(0, 0, 1),
		SourceFile:             filename,
		ParsedAt:               time.Now(),
		RawData:                nil, // Would be set from original JSON
	}
}

// parseDate parses a date string in ISO 8601 format.
func parseDate(dateStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, dateStr)
	return t
}

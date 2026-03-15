package cielo

import (
	"fmt"
	"io"
	"time"

	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// Parser implements the common.Parser interface for Cielo EEFI/EEVC format.
// Cielo provides settlement files in proprietary EEFI (Extrato Eletrônico de Financeira) and
// EEVC (Extrato Eletrônico de Vendas por Cartão) formats.
type Parser struct {
	Format string // "eefi" or "eevc"
}

// NewParser creates a new Cielo parser with the specified format.
func NewParser(format string) *Parser {
	return &Parser{
		Format: format,
	}
}

// Parse reads a Cielo settlement file and returns a slice of TransactionRecord.
func (p *Parser) Parse(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read Cielo file: %w", err)
	}

	var records []common.TransactionRecord

	// TODO: Implement Cielo EEFI/EEVC parsing
	// EEFI format typically includes settlement details with fixed-width fields
	// EEVC format contains card sales details
	// Both formats require bank-specific field position mapping

	// Stub: Parse the data and return records
	if len(data) > 0 {
		record := common.TransactionRecord{
			ID:                     "cielo-placeholder",
			SourceType:             common.SourceTypeCardCredit,
			ExternalID:             "CIELO001",
			CounterpartyDocument:   "00000000000000",
			AmountCents:            0,
			NetAmountCents:         0,
			TransactionDate:        time.Now(),
			ExpectedSettlementDate: time.Now().AddDate(0, 0, 1),
			SourceFile:             filename,
			ParsedAt:               time.Now(),
			RawData:                data,
		}
		records = append(records, record)
	}

	return records, nil
}

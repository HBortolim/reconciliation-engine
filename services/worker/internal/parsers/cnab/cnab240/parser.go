package cnab240

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hbortolim/reconciliation-engine/internal/parsers/cnab/profiles"
	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// Parser implements the common.Parser interface for CNAB 240 format.
// CNAB 240 uses fixed-width positional parsing with specific field alignments.
type Parser struct {
	BankProfile *profiles.BankProfile
}

// NewParser creates a new CNAB 240 parser with the given bank profile.
func NewParser(bankProfile *profiles.BankProfile) *Parser {
	return &Parser{
		BankProfile: bankProfile,
	}
}

// Parse reads a CNAB 240 file and returns a slice of TransactionRecord.
// CNAB 240 structure: Header (240 chars) + multiple detail records (240 chars each) + Trailer
func (p *Parser) Parse(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read CNAB 240 file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var records []common.TransactionRecord

	// CNAB 240 file structure:
	// Line 0: Header (record type 0)
	// Lines 1 to n-2: Detail records (record type 1 for segment A, 3 for segment J, etc)
	// Line n-1: Trailer (record type 9)

	for i := 1; i < len(lines)-1; i++ {
		line := strings.TrimSpace(lines[i])
		if len(line) < 240 {
			continue // Skip lines that don't meet CNAB 240 minimum length
		}

		record := p.parseDetailRecord(line, filename)
		if record != nil {
			records = append(records, *record)
		}
	}

	return records, nil
}

// parseDetailRecord extracts a TransactionRecord from a CNAB 240 detail line.
func (p *Parser) parseDetailRecord(line string, filename string) *common.TransactionRecord {
	if len(line) < 240 {
		return nil
	}

	// Field positions are bank-specific; using defaults for illustration
	// These would be overridden based on BankProfile settings

	// Extract key fields based on fixed positions
	// Positions are 1-indexed in CNAB documentation, but 0-indexed in code
	externalID := strings.TrimSpace(line[32:48]) // Nos. Sequencial do Registro (16 chars)
	amountStr := strings.TrimSpace(line[82:100]) // Valor (18 chars, rightmost 2 are decimals)
	dateStr := strings.TrimSpace(line[105:113])  // Data do Lançamento (8 chars, DDMMYYYY)

	var amountCents int64
	if len(amountStr) > 0 {
		fmt.Sscanf(amountStr, "%d", &amountCents)
	}

	transDate := parseDate(dateStr, "02012006") // DDMMYYYY format
	if transDate.IsZero() {
		transDate = time.Now()
	}

	return &common.TransactionRecord{
		ID:                     fmt.Sprintf("cnab240-%s", externalID),
		SourceType:             common.SourceTypeTED, // Would be determined from segment type
		ExternalID:             externalID,
		CounterpartyDocument:   strings.TrimSpace(line[54:68]), // CPF/CNPJ field
		AmountCents:            amountCents,
		NetAmountCents:         amountCents,
		TransactionDate:        transDate,
		ExpectedSettlementDate: transDate.AddDate(0, 0, 1),
		SourceFile:             filename,
		ParsedAt:               time.Now(),
		RawData:                []byte(line),
	}
}

// parseDate parses a date string in the given format.
func parseDate(dateStr, format string) time.Time {
	if len(dateStr) == 0 {
		return time.Time{}
	}
	t, _ := time.Parse(format, dateStr)
	return t
}

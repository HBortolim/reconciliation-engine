package cnab400

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hbortolim/reconciliation-engine/internal/parsers/cnab/profiles"
	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// Parser implements the common.Parser interface for CNAB 400 format.
// CNAB 400 is an older format with fixed-width records of 400 characters.
type Parser struct {
	BankProfile *profiles.BankProfile
}

// NewParser creates a new CNAB 400 parser with the given bank profile.
func NewParser(bankProfile *profiles.BankProfile) *Parser {
	return &Parser{
		BankProfile: bankProfile,
	}
}

// Parse reads a CNAB 400 file and returns a slice of TransactionRecord.
// CNAB 400 structure: Header (record type 0) + multiple detail records (record type 1) + Trailer (record type 9)
func (p *Parser) Parse(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read CNAB 400 file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var records []common.TransactionRecord

	// CNAB 400 file structure:
	// Line 0: Header (record type 0)
	// Lines 1 to n-2: Detail records (record type 1)
	// Line n-1: Trailer (record type 9)

	for i := 1; i < len(lines)-1; i++ {
		line := strings.TrimSpace(lines[i])
		if len(line) < 400 {
			continue // Skip lines that don't meet CNAB 400 minimum length
		}

		record := p.parseDetailRecord(line, filename)
		if record != nil {
			records = append(records, *record)
		}
	}

	return records, nil
}

// parseDetailRecord extracts a TransactionRecord from a CNAB 400 detail line.
func (p *Parser) parseDetailRecord(line string, filename string) *common.TransactionRecord {
	if len(line) < 400 {
		return nil
	}

	// Field positions for CNAB 400 (1-indexed in docs, 0-indexed in code)
	externalID := strings.TrimSpace(line[53:73])        // Número de inscrição (20 chars)
	amountStr := strings.TrimSpace(line[82:100])        // Valor (18 chars)
	dateStr := strings.TrimSpace(line[108:116])         // Data lançamento (8 chars, DDMMYYYY)
	counterpartyDoc := strings.TrimSpace(line[104:118]) // CPF/CNPJ

	var amountCents int64
	if len(amountStr) > 0 {
		fmt.Sscanf(amountStr, "%d", &amountCents)
	}

	transDate := parseDate(dateStr, "02012006") // DDMMYYYY format
	if transDate.IsZero() {
		transDate = time.Now()
	}

	return &common.TransactionRecord{
		ID:                     fmt.Sprintf("cnab400-%s", externalID),
		SourceType:             common.SourceTypeDOC,
		ExternalID:             externalID,
		CounterpartyDocument:   counterpartyDoc,
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

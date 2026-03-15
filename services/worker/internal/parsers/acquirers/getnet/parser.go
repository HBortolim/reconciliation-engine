package getnet

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// Parser implements the common.Parser interface for Getnet positional format.
// Getnet provides settlement files with fixed-width positional fields.
type Parser struct{}

// NewParser creates a new Getnet parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads a Getnet settlement file and returns a slice of TransactionRecord.
func (p *Parser) Parse(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read Getnet file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var records []common.TransactionRecord

	// Getnet files typically have:
	// Header line
	// Detail lines (fixed-width)
	// Footer line

	for i := 1; i < len(lines)-1; i++ {
		line := strings.TrimSpace(lines[i])
		if len(line) < 100 { // Minimum expected length for Getnet records
			continue
		}

		record := p.parsePositionalRecord(line, filename)
		if record != nil {
			records = append(records, *record)
		}
	}

	return records, nil
}

// parsePositionalRecord extracts a TransactionRecord from a fixed-width line.
func (p *Parser) parsePositionalRecord(line string, filename string) *common.TransactionRecord {
	if len(line) < 100 {
		return nil
	}

	// Extract fields based on Getnet positional format (positions vary by bank)
	// These are placeholder positions
	nsu := strings.TrimSpace(line[0:12])              // NSU (12 chars)
	amountStr := strings.TrimSpace(line[12:30])       // Amount (18 chars)
	dateStr := strings.TrimSpace(line[30:38])         // Date (8 chars, DDMMYYYY)
	counterpartyDoc := strings.TrimSpace(line[38:52]) // CNPJ/CPF (14 chars)

	if nsu == "" {
		return nil
	}

	var amountCents int64
	fmt.Sscanf(amountStr, "%d", &amountCents)

	transDate := parseDate(dateStr, "02012006") // DDMMYYYY format
	if transDate.IsZero() {
		transDate = time.Now()
	}

	return &common.TransactionRecord{
		ID:                     fmt.Sprintf("getnet-%s", nsu),
		SourceType:             common.SourceTypeCardCredit,
		NSU:                    nsu,
		ExternalID:             nsu,
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

package pagseguro

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// Parser implements the common.Parser interface for PagSeguro CSV format.
// PagSeguro provides settlement reports as CSV files with specific column structure.
type Parser struct{}

// NewParser creates a new PagSeguro parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads a PagSeguro CSV settlement file and returns a slice of TransactionRecord.
func (p *Parser) Parse(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ','

	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Create a map of column names to indices
	columnMap := make(map[string]int)
	for i, header := range headers {
		columnMap[strings.ToLower(strings.TrimSpace(header))] = i
	}

	var records []common.TransactionRecord

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		record := p.parseCSVRow(row, columnMap, filename)
		if record != nil {
			records = append(records, *record)
		}
	}

	return records, nil
}

// parseCSVRow parses a single CSV row into a TransactionRecord.
func (p *Parser) parseCSVRow(row []string, columnMap map[string]int, filename string) *common.TransactionRecord {
	// Extract fields based on typical PagSeguro CSV format
	referenceID := getColumnValue(row, columnMap, "reference")
	amountStr := getColumnValue(row, columnMap, "amount")
	dateStr := getColumnValue(row, columnMap, "date")
	buyerDoc := getColumnValue(row, columnMap, "buyer_document")

	if referenceID == "" {
		return nil // No reference ID means invalid transaction
	}

	var amountCents int64
	fmt.Sscanf(amountStr, "%d", &amountCents)

	transDate := parseDate(dateStr)
	if transDate.IsZero() {
		transDate = time.Now()
	}

	return &common.TransactionRecord{
		ID:                     fmt.Sprintf("pagseguro-%s", referenceID),
		SourceType:             common.SourceTypeCardCredit,
		ExternalID:             referenceID,
		CounterpartyDocument:   buyerDoc,
		AmountCents:            amountCents,
		NetAmountCents:         amountCents,
		TransactionDate:        transDate,
		ExpectedSettlementDate: transDate.AddDate(0, 0, 1),
		SourceFile:             filename,
		ParsedAt:               time.Now(),
		RawData:                []byte(strings.Join(row, "|")),
	}
}

// getColumnValue safely retrieves a value from a row by column name.
func getColumnValue(row []string, columnMap map[string]int, columnName string) string {
	idx, ok := columnMap[columnName]
	if !ok || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

// parseDate parses a date string in common Brazilian formats.
func parseDate(dateStr string) time.Time {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}
	}

	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"02-01-2006",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	return time.Time{}
}

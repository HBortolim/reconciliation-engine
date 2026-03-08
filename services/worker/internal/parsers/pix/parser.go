package pix

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// Parser implements the common.Parser interface for Pix settlement reports.
// Pix reports can be in JSON or CSV format.
type Parser struct {
	Format string // "json" or "csv"
}

// NewParser creates a new Pix parser with the specified format.
func NewParser(format string) *Parser {
	return &Parser{
		Format: strings.ToLower(format),
	}
}

// Parse reads a Pix settlement report and returns a slice of TransactionRecord.
func (p *Parser) Parse(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	if p.Format == "json" {
		return p.parseJSON(reader, filename)
	} else if p.Format == "csv" {
		return p.parseCSV(reader, filename)
	}
	return nil, fmt.Errorf("unsupported Pix format: %s", p.Format)
}

// parseJSON parses a JSON Pix settlement report.
func (p *Parser) parseJSON(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	var data interface{}
	if err := json.NewDecoder(reader).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var records []common.TransactionRecord
	// TODO: Parse JSON structure based on Pix settlement report format
	// Expected structure includes: e2eId, amount, settlementDate, counterpartyDocument, etc.

	return records, nil
}

// parseCSV parses a CSV Pix settlement report.
func (p *Parser) parseCSV(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';' // Brazilian CSV often uses semicolon

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
	// Extract fields based on column map
	e2eID := getColumnValue(row, columnMap, "e2eid")
	amountStr := getColumnValue(row, columnMap, "amount")
	dateStr := getColumnValue(row, columnMap, "settlementdate")
	counterpartyDoc := getColumnValue(row, columnMap, "counterpartydocument")

	if !validateE2EID(e2eID) {
		return nil // Invalid E2EID
	}

	var amountCentavos int64
	fmt.Sscanf(amountStr, "%d", &amountCentavos)

	transDate := parseDate(dateStr)
	if transDate.IsZero() {
		transDate = time.Now()
	}

	return &common.TransactionRecord{
		ID:                     e2eID,
		SourceType:             common.SourceTypePIX,
		E2EID:                  e2eID,
		ExternalID:             e2eID,
		CounterpartyDocument:   counterpartyDoc,
		AmountCentavos:         amountCentavos,
		NetAmountCentavos:      amountCentavos,
		TransactionDate:        transDate,
		ExpectedSettlementDate: transDate,
		SourceFile:             filename,
		ParsedAt:               time.Now(),
		RawData:                json.RawMessage(strings.Join(row, "|")),
	}
}

// validateE2EID validates a Pix E2EID using the standard format.
// E2EID format: AAAABBBBCCCCDDDD0XXXXX (35 chars total)
func validateE2EID(e2eID string) bool {
	// Pix E2EID validation regex (simplified)
	pattern := `^[A-Z0-9]{35}$`
	matched, _ := regexp.MatchString(pattern, e2eID)
	return matched
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

	// Try common Brazilian date formats
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

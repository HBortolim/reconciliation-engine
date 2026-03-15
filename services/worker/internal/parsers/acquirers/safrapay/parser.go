package safrapay

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// Parser implements the common.Parser interface for SafraPay format.
// SafraPay provides settlement data in CSV or fixed-width format.
type Parser struct {
	Format string // "csv" or "fixed"
}

// NewParser creates a new SafraPay parser with the specified format.
func NewParser(format string) *Parser {
	return &Parser{
		Format: strings.ToLower(format),
	}
}

// Parse reads a SafraPay settlement file and returns a slice of TransactionRecord.
func (p *Parser) Parse(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	if p.Format == "csv" {
		return p.parseCSV(reader, filename)
	}
	return p.parseFixed(reader, filename)
}

// parseCSV parses a SafraPay CSV settlement file.
func (p *Parser) parseCSV(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

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
	externalID := getColumnValue(row, columnMap, "id")
	amountStr := getColumnValue(row, columnMap, "amount")
	dateStr := getColumnValue(row, columnMap, "date")
	counterpartyDoc := getColumnValue(row, columnMap, "document")

	if externalID == "" {
		return nil
	}

	var amountCents int64
	fmt.Sscanf(amountStr, "%d", &amountCents)

	transDate := parseDate(dateStr)
	if transDate.IsZero() {
		transDate = time.Now()
	}

	return &common.TransactionRecord{
		ID:                     fmt.Sprintf("safrapay-%s", externalID),
		SourceType:             common.SourceTypeCardCredit,
		ExternalID:             externalID,
		CounterpartyDocument:   counterpartyDoc,
		AmountCents:            amountCents,
		NetAmountCents:         amountCents,
		TransactionDate:        transDate,
		ExpectedSettlementDate: transDate.AddDate(0, 0, 1),
		SourceFile:             filename,
		ParsedAt:               time.Now(),
		RawData:                []byte(strings.Join(row, "|")),
	}
}

// parseFixed parses a SafraPay fixed-width settlement file.
func (p *Parser) parseFixed(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read SafraPay file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var records []common.TransactionRecord

	for i := 1; i < len(lines)-1; i++ {
		line := strings.TrimSpace(lines[i])
		if len(line) < 80 {
			continue
		}

		record := p.parseFixedRecord(line, filename)
		if record != nil {
			records = append(records, *record)
		}
	}

	return records, nil
}

// parseFixedRecord extracts a TransactionRecord from a fixed-width line.
func (p *Parser) parseFixedRecord(line string, filename string) *common.TransactionRecord {
	if len(line) < 80 {
		return nil
	}

	externalID := strings.TrimSpace(line[0:20])       // ID (20 chars)
	amountStr := strings.TrimSpace(line[20:38])       // Amount (18 chars)
	dateStr := strings.TrimSpace(line[38:46])         // Date (8 chars)
	counterpartyDoc := strings.TrimSpace(line[46:60]) // Document (14 chars)

	if externalID == "" {
		return nil
	}

	var amountCents int64
	fmt.Sscanf(amountStr, "%d", &amountCents)

	transDate := parseDate(dateStr)
	if transDate.IsZero() {
		transDate = time.Now()
	}

	return &common.TransactionRecord{
		ID:                     fmt.Sprintf("safrapay-%s", externalID),
		SourceType:             common.SourceTypeCardCredit,
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

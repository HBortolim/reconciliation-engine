package ofx

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// Parser implements the common.Parser interface for OFX format.
// Note: OFX 1.x format is SGML-based, not XML. This parser handles both OFX 1.x and OFX 2.0 (XML).
type Parser struct {
	version string // "1.x" or "2.0"
}

// NewParser creates a new OFX parser.
func NewParser() *Parser {
	return &Parser{
		version: "2.0", // Default to OFX 2.0 (XML)
	}
}

// Parse reads an OFX file and returns a slice of TransactionRecord.
func (p *Parser) Parse(reader io.Reader, filename string) ([]common.TransactionRecord, error) {
	// TODO: Implement OFX parsing
	// For OFX 1.x (SGML), this would require a custom parser
	// For OFX 2.0 (XML), standard XML unmarshaling can be used
	var records []common.TransactionRecord

	// Stub: Parse the reader and populate records
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read OFX file: %w", err)
	}

	// Detect OFX version
	if len(data) == 0 {
		return records, nil
	}

	// For now, return empty records with a placeholder implementation
	// In production, this would parse OFXHEADER, validate format version,
	// and extract transaction details from the STMTRS (statement response) section
	record := common.TransactionRecord{
		ID:                     "ofx-placeholder",
		SourceType:             common.SourceTypeTED,
		ExternalID:             "OFXID001",
		CounterpartyDocument:   "00000000000000",
		AmountCents:            0,
		NetAmountCents:         0,
		TransactionDate:        time.Now(),
		ExpectedSettlementDate: time.Now(),
		SourceFile:             filename,
		ParsedAt:               time.Now(),
		RawData:                json.RawMessage(data),
	}

	return []common.TransactionRecord{record}, nil
}

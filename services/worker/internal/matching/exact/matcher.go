package exact

import (
	"time"

	"github.com/hbortolim/reconciliation-engine/internal/matching"
	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// Matcher implements exact matching by comparing key identifiers.
// It matches records based on NSU, E2EID, or NossoNumero depending on source type.
type Matcher struct{}

// NewMatcher creates a new ExactMatcher.
func NewMatcher() *Matcher {
	return &Matcher{}
}

// Match returns matched pairs, unmatched expected, and unmatched actual records.
func (m *Matcher) Match(expected, actual []common.TransactionRecord) (
	[]matching.ReconciliationPair,
	[]common.TransactionRecord,
	[]common.TransactionRecord) {

	var pairs []matching.ReconciliationPair
	matchedExpectedIndices := make(map[int]bool)
	matchedActualIndices := make(map[int]bool)

	// Try to match each expected record with an actual record
	for i, exp := range expected {
		for j, act := range actual {
			if matchedActualIndices[j] {
				continue // Already matched
			}

			if m.exactMatch(&exp, &act) {
				pairs = append(pairs, matching.ReconciliationPair{
					Expected:        exp,
					Actual:          act,
					MatchType:       matching.MatchTypeExact,
					ConfidenceScore: 1.0,
				})
				matchedExpectedIndices[i] = true
				matchedActualIndices[j] = true
				break
			}
		}
	}

	// Collect unmatched records
	var unmatchedExpected []common.TransactionRecord
	for i, exp := range expected {
		if !matchedExpectedIndices[i] {
			unmatchedExpected = append(unmatchedExpected, exp)
		}
	}

	var unmatchedActual []common.TransactionRecord
	for i, act := range actual {
		if !matchedActualIndices[i] {
			unmatchedActual = append(unmatchedActual, act)
		}
	}

	return pairs, unmatchedExpected, unmatchedActual
}

// exactMatch checks if two records match based on source-specific identifiers.
func (m *Matcher) exactMatch(expected, actual *common.TransactionRecord) bool {
	// Check by source-specific identifier
	switch expected.SourceType {
	case common.SourceTypePIX:
		return expected.E2EID != "" && expected.E2EID == actual.E2EID
	case common.SourceTypeCardCredit, common.SourceTypeCardDebit:
		return expected.NSU != "" && expected.NSU == actual.NSU
	case common.SourceTypeBOLETO, common.SourceTypeDOC, common.SourceTypeTED:
		return expected.NossoNumero != "" && expected.NossoNumero == actual.NossoNumero
	default:
		// Fallback: match by external ID and amount
		return expected.ExternalID == actual.ExternalID &&
			expected.AmountCents == actual.AmountCents &&
			m.datesAreClose(expected.TransactionDate, actual.TransactionDate)
	}
}

// datesAreClose checks if two dates are within 1 day of each other.
func (m *Matcher) datesAreClose(date1, date2 time.Time) bool {
	diff := date1.Sub(date2)
	if diff < 0 {
		diff = -diff
	}
	return diff.Hours() <= 24
}

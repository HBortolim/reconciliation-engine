package fuzzy

import (
	"github.com/hbortolim/reconciliation-engine/internal/matching"
	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// MatcherConfig contains configuration for fuzzy matching.
type MatcherConfig struct {
	AmountToleranceCents  int64
	DateWindowDays        int
	CounterpartyThreshold float64
}

// DefaultConfig returns default fuzzy matcher configuration.
func DefaultConfig() *MatcherConfig {
	return &MatcherConfig{
		AmountToleranceCents:  100,  // 1 BRL tolerance
		DateWindowDays:        3,    // 3-day window
		CounterpartyThreshold: 0.80, // 80% similarity
	}
}

// Matcher implements fuzzy matching with configurable thresholds.
// It scores matches based on amount tolerance, date window, and counterparty name similarity.
type Matcher struct {
	Config *MatcherConfig
}

// NewMatcher creates a new FuzzyMatcher with default configuration.
func NewMatcher() *Matcher {
	return &Matcher{
		Config: DefaultConfig(),
	}
}

// NewMatcherWithConfig creates a new FuzzyMatcher with custom configuration.
func NewMatcherWithConfig(config *MatcherConfig) *Matcher {
	return &Matcher{
		Config: config,
	}
}

// Match returns matched pairs with confidence scores, plus unmatched records.
func (m *Matcher) Match(expected, actual []common.TransactionRecord) (
	[]matching.ReconciliationPair,
	[]common.TransactionRecord,
	[]common.TransactionRecord) {

	var pairs []matching.ReconciliationPair
	matchedExpectedIndices := make(map[int]bool)
	matchedActualIndices := make(map[int]bool)

	// Try to find the best match for each expected record
	for i, exp := range expected {
		bestScore := 0.0
		bestActualIdx := -1

		for j, act := range actual {
			if matchedActualIndices[j] {
				continue
			}

			score := m.scoreMatch(&exp, &act)
			if score > bestScore && score > 0.5 { // Minimum score threshold
				bestScore = score
				bestActualIdx = j
			}
		}

		if bestActualIdx >= 0 {
			act := actual[bestActualIdx]
			pairs = append(pairs, matching.ReconciliationPair{
				Expected:           exp,
				Actual:             act,
				MatchType:          matching.MatchTypeFuzzy,
				ConfidenceScore:    bestScore,
				DiscrepancyDetails: m.computeDiscrepancies(&exp, &act),
			})
			matchedExpectedIndices[i] = true
			matchedActualIndices[bestActualIdx] = true
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

// scoreMatch computes a match score between 0 and 1 for two records.
func (m *Matcher) scoreMatch(expected, actual *common.TransactionRecord) float64 {
	// Amount score
	amountDiff := expected.AmountCents - actual.AmountCents
	if amountDiff < 0 {
		amountDiff = -amountDiff
	}

	amountScore := 1.0
	if amountDiff > m.Config.AmountToleranceCents {
		if amountDiff > m.Config.AmountToleranceCents*2 {
			return 0 // Disqualify if amount difference is too large
		}
		amountScore = 1.0 - (float64(amountDiff) / float64(m.Config.AmountToleranceCents*2))
	}

	// Date score
	dateDiff := expected.TransactionDate.Sub(actual.TransactionDate)
	if dateDiff < 0 {
		dateDiff = -dateDiff
	}

	dateScore := 1.0
	daysDiff := int(dateDiff.Hours() / 24)
	if daysDiff > m.Config.DateWindowDays {
		return 0 // Disqualify if dates are too far apart
	}
	if daysDiff > 0 {
		dateScore = 1.0 - (float64(daysDiff) / float64(m.Config.DateWindowDays))
	}

	// Counterparty score (if names are available)
	counterpartyScore := 1.0
	if expected.CounterpartyName != "" && actual.CounterpartyName != "" {
		counterpartyScore = SimilarityScore(expected.CounterpartyName, actual.CounterpartyName)
	}

	// Weighted average
	return (amountScore * 0.4) + (dateScore * 0.3) + (counterpartyScore * 0.3)
}

// computeDiscrepancies identifies differences between two matched records.
func (m *Matcher) computeDiscrepancies(expected, actual *common.TransactionRecord) *matching.DiscrepancyDetails {
	details := &matching.DiscrepancyDetails{
		AmountDifferenceCents: expected.AmountCents - actual.AmountCents,
		FeeDifferenceCents:    expected.FeeCents - actual.FeeCents,
	}

	dateDiff := expected.TransactionDate.Sub(actual.TransactionDate)
	details.DateDifferenceDays = int(dateDiff.Hours() / 24)

	if expected.CounterpartyDocument != actual.CounterpartyDocument {
		details.CounterpartyMismatch = true
	}

	return details
}

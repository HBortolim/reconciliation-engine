package matching

import "github.com/hbortolim/reconciliation-engine/internal/parsers/common"

// MatchType represents the type of reconciliation match.
type MatchType string

const (
	MatchTypeExact     MatchType = "EXACT"
	MatchTypeFuzzy     MatchType = "FUZZY"
	MatchTypeAggregate MatchType = "AGGREGATE"
)

// DiscrepancyDetails contains information about differences between matched records.
type DiscrepancyDetails struct {
	AmountDifferenceCents int64  `json:"amount_difference_cents"`
	DateDifferenceDays    int    `json:"date_difference_days"`
	CounterpartyMismatch  bool   `json:"counterparty_mismatch"`
	FeeDifferenceCents    int64  `json:"fee_difference_cents"`
	Notes                 string `json:"notes"`
}

// ReconciliationPair represents a matched pair of expected and actual transactions.
type ReconciliationPair struct {
	Expected           common.TransactionRecord `json:"expected"`
	Actual             common.TransactionRecord `json:"actual"`
	MatchType          MatchType                `json:"match_type"`
	ConfidenceScore    float64                  `json:"confidence_score"`
	DiscrepancyDetails *DiscrepancyDetails      `json:"discrepancy_details,omitempty"`
}

// ReconciliationResult contains the results of a reconciliation run.
type ReconciliationResult struct {
	Pairs              []ReconciliationPair       `json:"pairs"`
	UnmatchedExpected  []common.TransactionRecord `json:"unmatched_expected"`
	UnmatchedActual    []common.TransactionRecord `json:"unmatched_actual"`
	TotalExpected      int                        `json:"total_expected"`
	TotalActual        int                        `json:"total_actual"`
	MatchedCount       int                        `json:"matched_count"`
	ReconciliationRate float64                    `json:"reconciliation_rate"`
}

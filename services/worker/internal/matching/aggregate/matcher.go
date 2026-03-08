package aggregate

import (
	"fmt"
	"sort"

	"github.com/hbortolim/reconciliation-engine/internal/matching"
	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// MatcherConfig contains configuration for aggregate matching.
type MatcherConfig struct {
	AmountToleranceCentavos int64
	MaxSubsetSize           int // Maximum number of transactions to group (default 50)
}

// DefaultConfig returns default aggregate matcher configuration.
func DefaultConfig() *MatcherConfig {
	return &MatcherConfig{
		AmountToleranceCentavos: 100,
		MaxSubsetSize:           50,
	}
}

// Matcher implements N-to-1 batch settlement matching using a branch-and-bound algorithm.
type Matcher struct {
	Config *MatcherConfig
}

// NewMatcher creates a new AggregateMatcher with default configuration.
func NewMatcher() *Matcher {
	return &Matcher{
		Config: DefaultConfig(),
	}
}

// NewMatcherWithConfig creates a new AggregateMatcher with custom configuration.
func NewMatcherWithConfig(config *MatcherConfig) *Matcher {
	return &Matcher{
		Config: config,
	}
}

// Match finds aggregated matches where multiple expected transactions settle as one.
func (m *Matcher) Match(expected, actual []common.TransactionRecord) (
	[]matching.ReconciliationPair,
	[]common.TransactionRecord,
	[]common.TransactionRecord) {

	var pairs []matching.ReconciliationPair
	matchedExpectedIndices := make(map[int]bool)
	matchedActualIndices := make(map[int]bool)

	// For each actual record, try to find a subset of expected records that sum to it
	for j, act := range actual {
		if matchedActualIndices[j] {
			continue
		}

		// Find best subset match using branch-and-bound
		subset := m.findBestSubset(expected, &act, matchedExpectedIndices)
		if subset != nil && len(subset) > 0 {
			// Create aggregated pair
			pair := matching.ReconciliationPair{
				Expected:        subset[0], // Use first as representative
				Actual:          act,
				MatchType:       matching.MatchTypeAggregate,
				ConfidenceScore: m.scoreSubsetMatch(subset, &act),
			}
			pairs = append(pairs, pair)

			for _, exp := range subset {
				// Find and mark index
				for i, e := range expected {
					if !matchedExpectedIndices[i] && e.ID == exp.ID {
						matchedExpectedIndices[i] = true
						break
					}
				}
			}
			matchedActualIndices[j] = true
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

// findBestSubset uses branch-and-bound to find a subset of expected records that best match actual.
func (m *Matcher) findBestSubset(expected []common.TransactionRecord, actual *common.TransactionRecord,
	alreadyMatched map[int]bool) []common.TransactionRecord {

	// Filter available candidates
	var candidates []common.TransactionRecord
	for i, exp := range expected {
		if !alreadyMatched[i] {
			candidates = append(candidates, exp)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	// Limit search space to MaxSubsetSize
	if len(candidates) > m.Config.MaxSubsetSize {
		// Sort by amount proximity and take top candidates
		sort.Slice(candidates, func(i, j int) bool {
			diffI := actual.AmountCentavos - candidates[i].AmountCentavos
			if diffI < 0 {
				diffI = -diffI
			}
			diffJ := actual.AmountCentavos - candidates[j].AmountCentavos
			if diffJ < 0 {
				diffJ = -diffJ
			}
			return diffI < diffJ
		})
		candidates = candidates[:m.Config.MaxSubsetSize]
	}

	// Use dynamic programming to find best subset sum
	return m.solveSubsetSum(candidates, actual.AmountCentavos)
}

// solveSubsetSum finds a subset of transactions that sum closest to the target amount.
func (m *Matcher) solveSubsetSum(transactions []common.TransactionRecord, targetAmount int64) []common.TransactionRecord {
	n := len(transactions)
	if n == 0 {
		return nil
	}

	// DP table: dp[i][j] = can we achieve sum j using first i items?
	maxSum := targetAmount + (m.Config.AmountToleranceCentavos * 2)
	if maxSum < 0 {
		maxSum = targetAmount + 10000
	}

	// For efficiency, limit the DP table size
	if maxSum > 1000000 {
		maxSum = 1000000
	}

	dp := make([][]bool, n+1)
	parent := make([][]int, n+1)

	for i := range dp {
		dp[i] = make([]bool, maxSum+1)
		parent[i] = make([]int, maxSum+1)
		for j := range dp[i] {
			parent[i][j] = -1
		}
	}

	dp[0][0] = true

	// Fill DP table
	for i := 1; i <= n; i++ {
		amount := transactions[i-1].AmountCentavos
		for j := 0; j <= maxSum; j++ {
			// Don't take item i-1
			if dp[i-1][j] {
				dp[i][j] = true
				parent[i][j] = 0
			}

			// Take item i-1
			if j >= amount && dp[i-1][j-amount] {
				dp[i][j] = true
				parent[i][j] = 1
			}
		}
	}

	// Find the sum closest to target
	bestSum := 0
	for j := 0; j <= maxSum; j++ {
		if dp[n][j] {
			diff := targetAmount - int64(j)
			if diff < 0 {
				diff = -diff
			}
			bestDiff := targetAmount - int64(bestSum)
			if bestDiff < 0 {
				bestDiff = -bestDiff
			}

			if diff < bestDiff || bestSum == 0 {
				bestSum = j
			}
		}
	}

	if bestSum == 0 {
		return nil
	}

	// Reconstruct subset
	var result []common.TransactionRecord
	j := bestSum
	for i := n; i > 0 && j > 0; i-- {
		if parent[i][j] == 1 {
			result = append([]common.TransactionRecord{transactions[i-1]}, result...)
			j -= int(transactions[i-1].AmountCentavos)
		}
	}

	return result
}

// scoreSubsetMatch computes a confidence score for an aggregate match.
func (m *Matcher) scoreSubsetMatch(subset []common.TransactionRecord, actual *common.TransactionRecord) float64 {
	totalAmount := int64(0)
	for _, txn := range subset {
		totalAmount += txn.AmountCentavos
	}

	diff := totalAmount - actual.AmountCentavos
	if diff < 0 {
		diff = -diff
	}

	if diff > m.Config.AmountToleranceCentavos*2 {
		return 0
	}

	return 1.0 - (float64(diff) / float64(m.Config.AmountToleranceCentavos*2))
}

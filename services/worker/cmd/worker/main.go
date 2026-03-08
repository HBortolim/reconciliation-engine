package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hbortolim/reconciliation-engine/internal/matching/exact"
	"github.com/hbortolim/reconciliation-engine/internal/matching/fuzzy"
	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

// Job represents a reconciliation job from the queue.
type Job struct {
	ID            string                        `json:"id"`
	ExpectedFile  string                        `json:"expected_file"`
	ActualFile    string                        `json:"actual_file"`
	SourceType    string                        `json:"source_type"`
	MatchStrategy string                        `json:"match_strategy"` // "exact" or "fuzzy"
	Timestamp     time.Time                     `json:"timestamp"`
}

// JobResult represents the result of a reconciliation job.
type JobResult struct {
	JobID             string    `json:"job_id"`
	Status            string    `json:"status"` // "success" or "error"
	MatchedCount      int       `json:"matched_count"`
	UnmatchedExpected int       `json:"unmatched_expected"`
	UnmatchedActual   int       `json:"unmatched_actual"`
	Error             string    `json:"error,omitempty"`
	CompletedAt       time.Time `json:"completed_at"`
}

func main() {
	log.Println("Reconciliation Engine Worker starting...")

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Connect to Redis stream for job queue
	// For now, this demonstrates the worker structure

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Example: Listen for jobs (stub implementation)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				log.Println("Waiting for jobs...")
				// TODO: Poll Redis stream for new jobs
			}
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received, gracefully closing...")
	cancel()
	log.Println("Worker stopped")
}

// processJob handles a single reconciliation job.
func processJob(job *Job) *JobResult {
	result := &JobResult{
		JobID:       job.ID,
		CompletedAt: time.Now(),
	}

	log.Printf("Processing job: %s (strategy: %s)\n", job.ID, job.MatchStrategy)

	// TODO: Parse both files based on source type
	var expectedRecords []common.TransactionRecord
	var actualRecords []common.TransactionRecord

	// Select matching strategy
	var pairs int
	var unmatched int
	var unmatchedActual int

	if job.MatchStrategy == "fuzzy" {
		matcher := fuzzy.NewMatcher()
		matchedPairs, unmatchedExp, unmatchedAct := matcher.Match(expectedRecords, actualRecords)
		pairs = len(matchedPairs)
		unmatched = len(unmatchedExp)
		unmatchedActual = len(unmatchedAct)
	} else {
		// Default to exact matching
		matcher := exact.NewMatcher()
		matchedPairs, unmatchedExp, unmatchedAct := matcher.Match(expectedRecords, actualRecords)
		pairs = len(matchedPairs)
		unmatched = len(unmatchedExp)
		unmatchedActual = len(unmatchedAct)
	}

	result.Status = "success"
	result.MatchedCount = pairs
	result.UnmatchedExpected = unmatched
	result.UnmatchedActual = unmatchedActual

	// TODO: Persist results to database
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	log.Printf("Job result: %s\n", string(resultJSON))

	return result
}

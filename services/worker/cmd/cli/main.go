package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hbortolim/reconciliation-engine/internal/matching/exact"
	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

func main() {
	// Define command-line flags
	filePath := flag.String("file", "", "Path to the input file for reconciliation")
	sourceType := flag.String("source-type", "", "Type of transaction source (PIX, BOLETO, CARD_CREDIT, etc.)")
	runID := flag.String("run-id", "", "Unique run ID for this reconciliation")
	debugMode := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	// Validate required flags
	if *filePath == "" {
		fmt.Println("Error: --file flag is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *sourceType == "" {
		fmt.Println("Error: --source-type flag is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *debugMode {
		log.Printf("Debug mode enabled\n")
		log.Printf("File: %s\n", *filePath)
		log.Printf("Source Type: %s\n", *sourceType)
		log.Printf("Run ID: %s\n", *runID)
	}

	// Read and parse the input file
	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// TODO: Select appropriate parser based on source type
	// For now, this is a stub that demonstrates the CLI structure

	log.Printf("Processing file: %s with source type: %s\n", *filePath, *sourceType)

	// Example: Create some dummy records and perform exact matching
	expectedRecords := []common.TransactionRecord{}
	actualRecords := []common.TransactionRecord{}

	matcher := exact.NewMatcher()
	pairs, unmatched, unmatchedActual := matcher.Match(expectedRecords, actualRecords)

	// Report results
	fmt.Printf("\n=== Reconciliation Results ===\n")
	fmt.Printf("Matched pairs: %d\n", len(pairs))
	fmt.Printf("Unmatched expected: %d\n", len(unmatched))
	fmt.Printf("Unmatched actual: %d\n", len(unmatchedActual))

	if *debugMode {
		for _, pair := range pairs {
			log.Printf("Match: %s <-> %s (confidence: %.2f)\n",
				pair.Expected.ID, pair.Actual.ID, pair.ConfidenceScore)
		}
	}
}

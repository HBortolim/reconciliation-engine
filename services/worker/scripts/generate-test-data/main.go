package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hbortolim/reconciliation-engine/internal/parsers/common"
)

func main() {
	fmt.Println("Test Data Generator for Reconciliation Engine")
	fmt.Println("============================================")
	fmt.Println()

	// Generate test data for each format
	generatePixTestData()
	generateCNABTestData()
	generateCardTestData()
}

// generatePixTestData creates synthetic Pix transaction records.
func generatePixTestData() {
	fmt.Println("Generating Pix test data...")
	records := make([]common.TransactionRecord, 10)

	for i := 0; i < 10; i++ {
		e2eID := fmt.Sprintf("20240308%s", uuid.New().String()[:20])
		records[i] = common.TransactionRecord{
			ID:                     e2eID,
			SourceType:             common.SourceTypePIX,
			E2EID:                  e2eID,
			ExternalID:             e2eID,
			CounterpartyDocument:   generateDocument(),
			CounterpartyName:       generateCounterpartyName(),
			AmountCentavos:         int64(rand.Intn(100000)),
			NetAmountCentavos:      int64(rand.Intn(100000)),
			TransactionDate:        time.Now().AddDate(0, 0, -rand.Intn(7)),
			ExpectedSettlementDate: time.Now(),
			SourceFile:             "pix-test.csv",
			ParsedAt:               time.Now(),
		}
	}

	// Write to JSON file
	data, _ := json.MarshalIndent(records, "", "  ")
	_ = os.WriteFile("test_data_pix.json", data, 0644)
	fmt.Printf("  Generated %d Pix records -> test_data_pix.json\n", len(records))
}

// generateCNABTestData creates synthetic CNAB 240 transaction records.
func generateCNABTestData() {
	fmt.Println("Generating CNAB 240 test data...")
	records := make([]common.TransactionRecord, 15)

	for i := 0; i < 15; i++ {
		nossoNumero := fmt.Sprintf("%015d", rand.Int63n(1000000000000000))
		records[i] = common.TransactionRecord{
			ID:                     fmt.Sprintf("cnab240-%d", i),
			SourceType:             common.SourceTypeTED,
			NossoNumero:            nossoNumero,
			ExternalID:             nossoNumero,
			CounterpartyDocument:   generateDocument(),
			CounterpartyName:       generateCounterpartyName(),
			AmountCentavos:         int64(rand.Intn(500000)),
			FeeCentavos:            int64(rand.Intn(1000)),
			NetAmountCentavos:      int64(rand.Intn(500000)),
			TransactionDate:        time.Now().AddDate(0, 0, -rand.Intn(5)),
			ExpectedSettlementDate: time.Now().AddDate(0, 0, 1),
			SourceFile:             "cnab240-test.txt",
			ParsedAt:               time.Now(),
		}
	}

	data, _ := json.MarshalIndent(records, "", "  ")
	_ = os.WriteFile("test_data_cnab240.json", data, 0644)
	fmt.Printf("  Generated %d CNAB 240 records -> test_data_cnab240.json\n", len(records))
}

// generateCardTestData creates synthetic card transaction records (Cielo, Rede, etc.).
func generateCardTestData() {
	fmt.Println("Generating Card transaction test data...")
	records := make([]common.TransactionRecord, 20)

	sources := []common.SourceType{
		common.SourceTypeCardCredit,
		common.SourceTypeCardDebit,
	}

	for i := 0; i < 20; i++ {
		nsu := fmt.Sprintf("%012d", rand.Int63n(1000000000000))
		sourceType := sources[rand.Intn(len(sources))]
		records[i] = common.TransactionRecord{
			ID:                     fmt.Sprintf("card-%d", i),
			SourceType:             sourceType,
			NSU:                    nsu,
			ExternalID:             nsu,
			CounterpartyDocument:   generateDocument(),
			CounterpartyName:       generateCounterpartyName(),
			AmountCentavos:         int64(rand.Intn(100000)),
			FeeCentavos:            int64(rand.Intn(500)),
			NetAmountCentavos:      int64(rand.Intn(100000)),
			TransactionDate:        time.Now().AddDate(0, 0, -rand.Intn(3)),
			ExpectedSettlementDate: time.Now().AddDate(0, 0, 1),
			SourceFile:             "card-test.csv",
			ParsedAt:               time.Now(),
		}
	}

	data, _ := json.MarshalIndent(records, "", "  ")
	_ = os.WriteFile("test_data_cards.json", data, 0644)
	fmt.Printf("  Generated %d Card records -> test_data_cards.json\n", len(records))
}

// generateDocument creates a random Brazilian document (CPF or CNPJ).
func generateDocument() string {
	if rand.Intn(2) == 0 {
		// Generate CPF (11 digits)
		return fmt.Sprintf("%011d", rand.Int63n(100000000000))
	} else {
		// Generate CNPJ (14 digits)
		return fmt.Sprintf("%014d", rand.Int63n(100000000000000))
	}
}

// generateCounterpartyName generates a random company or person name.
func generateCounterpartyName() string {
	companies := []string{
		"ABC Comercio LTDA",
		"XYZ Services S.A.",
		"Tech Solutions LTDA",
		"Global Enterprises Inc",
		"Local Market Ltda",
		"Innovation Systems S.A.",
		"Smart Business LTDA",
	}
	return companies[rand.Intn(len(companies))]
}

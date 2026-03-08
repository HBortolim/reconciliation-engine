package profiles

// FieldPosition defines the start and end position of a field in a fixed-width record.
type FieldPosition struct {
	Start int
	End   int
	Name  string
}

// BankProfile defines bank-specific field positions and rules for CNAB parsing.
type BankProfile struct {
	BankCode              string
	BankName              string
	Format                string // "240" or "400"
	ExternalIDPositions   FieldPosition
	AmountPositions       FieldPosition
	DatePositions         FieldPosition
	CounterpartyDocPositions FieldPosition
	NossoNumeroPositions  FieldPosition
}

// GetITAUProfile returns the Itaú bank profile.
func GetITAUProfile(format string) *BankProfile {
	return &BankProfile{
		BankCode:   "341",
		BankName:   "Itaú",
		Format:     format,
		ExternalIDPositions:   FieldPosition{Start: 32, End: 48, Name: "ExternalID"},
		AmountPositions:       FieldPosition{Start: 82, End: 100, Name: "Amount"},
		DatePositions:         FieldPosition{Start: 105, End: 113, Name: "Date"},
		CounterpartyDocPositions: FieldPosition{Start: 54, End: 68, Name: "CounterpartyDoc"},
		NossoNumeroPositions:  FieldPosition{Start: 100, End: 120, Name: "NossoNumero"},
	}
}

// GetBradescoProfile returns the Bradesco bank profile.
func GetBradescoProfile(format string) *BankProfile {
	return &BankProfile{
		BankCode:   "237",
		BankName:   "Bradesco",
		Format:     format,
		ExternalIDPositions:   FieldPosition{Start: 32, End: 48, Name: "ExternalID"},
		AmountPositions:       FieldPosition{Start: 82, End: 100, Name: "Amount"},
		DatePositions:         FieldPosition{Start: 105, End: 113, Name: "Date"},
		CounterpartyDocPositions: FieldPosition{Start: 54, End: 68, Name: "CounterpartyDoc"},
		NossoNumeroPositions:  FieldPosition{Start: 100, End: 120, Name: "NossoNumero"},
	}
}

// GetBBProfile returns the Banco do Brasil profile.
func GetBBProfile(format string) *BankProfile {
	return &BankProfile{
		BankCode:   "001",
		BankName:   "Banco do Brasil",
		Format:     format,
		ExternalIDPositions:   FieldPosition{Start: 32, End: 48, Name: "ExternalID"},
		AmountPositions:       FieldPosition{Start: 82, End: 100, Name: "Amount"},
		DatePositions:         FieldPosition{Start: 105, End: 113, Name: "Date"},
		CounterpartyDocPositions: FieldPosition{Start: 54, End: 68, Name: "CounterpartyDoc"},
		NossoNumeroPositions:  FieldPosition{Start: 100, End: 120, Name: "NossoNumero"},
	}
}

// GetSantanderProfile returns the Santander bank profile.
func GetSantanderProfile(format string) *BankProfile {
	return &BankProfile{
		BankCode:   "033",
		BankName:   "Santander",
		Format:     format,
		ExternalIDPositions:   FieldPosition{Start: 32, End: 48, Name: "ExternalID"},
		AmountPositions:       FieldPosition{Start: 82, End: 100, Name: "Amount"},
		DatePositions:         FieldPosition{Start: 105, End: 113, Name: "Date"},
		CounterpartyDocPositions: FieldPosition{Start: 54, End: 68, Name: "CounterpartyDoc"},
		NossoNumeroPositions:  FieldPosition{Start: 100, End: 120, Name: "NossoNumero"},
	}
}

// GetCaixaProfile returns the Caixa bank profile.
func GetCaixaProfile(format string) *BankProfile {
	return &BankProfile{
		BankCode:   "104",
		BankName:   "Caixa",
		Format:     format,
		ExternalIDPositions:   FieldPosition{Start: 32, End: 48, Name: "ExternalID"},
		AmountPositions:       FieldPosition{Start: 82, End: 100, Name: "Amount"},
		DatePositions:         FieldPosition{Start: 105, End: 113, Name: "Date"},
		CounterpartyDocPositions: FieldPosition{Start: 54, End: 68, Name: "CounterpartyDoc"},
		NossoNumeroPositions:  FieldPosition{Start: 100, End: 120, Name: "NossoNumero"},
	}
}

// GetBTGPactualProfile returns the BTG Pactual bank profile.
func GetBTGPactualProfile(format string) *BankProfile {
	return &BankProfile{
		BankCode:   "208",
		BankName:   "BTG Pactual",
		Format:     format,
		ExternalIDPositions:   FieldPosition{Start: 32, End: 48, Name: "ExternalID"},
		AmountPositions:       FieldPosition{Start: 82, End: 100, Name: "Amount"},
		DatePositions:         FieldPosition{Start: 105, End: 113, Name: "Date"},
		CounterpartyDocPositions: FieldPosition{Start: 54, End: 68, Name: "CounterpartyDoc"},
		NossoNumeroPositions:  FieldPosition{Start: 100, End: 120, Name: "NossoNumero"},
	}
}

// GetInterProfile returns the Inter bank profile.
func GetInterProfile(format string) *BankProfile {
	return &BankProfile{
		BankCode:   "077",
		BankName:   "Inter",
		Format:     format,
		ExternalIDPositions:   FieldPosition{Start: 32, End: 48, Name: "ExternalID"},
		AmountPositions:       FieldPosition{Start: 82, End: 100, Name: "Amount"},
		DatePositions:         FieldPosition{Start: 105, End: 113, Name: "Date"},
		CounterpartyDocPositions: FieldPosition{Start: 54, End: 68, Name: "CounterpartyDoc"},
		NossoNumeroPositions:  FieldPosition{Start: 100, End: 120, Name: "NossoNumero"},
	}
}

// GetNubankProfile returns the Nubank profile.
func GetNubankProfile(format string) *BankProfile {
	return &BankProfile{
		BankCode:   "260",
		BankName:   "Nubank",
		Format:     format,
		ExternalIDPositions:   FieldPosition{Start: 32, End: 48, Name: "ExternalID"},
		AmountPositions:       FieldPosition{Start: 82, End: 100, Name: "Amount"},
		DatePositions:         FieldPosition{Start: 105, End: 113, Name: "Date"},
		CounterpartyDocPositions: FieldPosition{Start: 54, End: 68, Name: "CounterpartyDoc"},
		NossoNumeroPositions:  FieldPosition{Start: 100, End: 120, Name: "NossoNumero"},
	}
}

// GetSicrediProfile returns the Sicredi bank profile.
func GetSicrediProfile(format string) *BankProfile {
	return &BankProfile{
		BankCode:   "748",
		BankName:   "Sicredi",
		Format:     format,
		ExternalIDPositions:   FieldPosition{Start: 32, End: 48, Name: "ExternalID"},
		AmountPositions:       FieldPosition{Start: 82, End: 100, Name: "Amount"},
		DatePositions:         FieldPosition{Start: 105, End: 113, Name: "Date"},
		CounterpartyDocPositions: FieldPosition{Start: 54, End: 68, Name: "CounterpartyDoc"},
		NossoNumeroPositions:  FieldPosition{Start: 100, End: 120, Name: "NossoNumero"},
	}
}

// GetSicoobProfile returns the Sicoob bank profile.
func GetSicoobProfile(format string) *BankProfile {
	return &BankProfile{
		BankCode:   "756",
		BankName:   "Sicoob",
		Format:     format,
		ExternalIDPositions:   FieldPosition{Start: 32, End: 48, Name: "ExternalID"},
		AmountPositions:       FieldPosition{Start: 82, End: 100, Name: "Amount"},
		DatePositions:         FieldPosition{Start: 105, End: 113, Name: "Date"},
		CounterpartyDocPositions: FieldPosition{Start: 54, End: 68, Name: "CounterpartyDoc"},
		NossoNumeroPositions:  FieldPosition{Start: 100, End: 120, Name: "NossoNumero"},
	}
}

// GetProfileByBankCode returns a bank profile by bank code and format.
func GetProfileByBankCode(bankCode, format string) *BankProfile {
	switch bankCode {
	case "341":
		return GetITAUProfile(format)
	case "237":
		return GetBradescoProfile(format)
	case "001":
		return GetBBProfile(format)
	case "033":
		return GetSantanderProfile(format)
	case "104":
		return GetCaixaProfile(format)
	case "208":
		return GetBTGPactualProfile(format)
	case "077":
		return GetInterProfile(format)
	case "260":
		return GetNubankProfile(format)
	case "748":
		return GetSicrediProfile(format)
	case "756":
		return GetSicoobProfile(format)
	default:
		return GetITAUProfile(format) // Default profile
	}
}

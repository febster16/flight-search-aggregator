package partners

// Provider names
const (
	ProviderGaruda  = "Garuda Indonesia"
	ProviderLion    = "Lion Air"
	ProviderBatik   = "Batik Air"
	ProviderAirAsia = "AirAsia"
)

// Currencies
const (
	CurrencyIDR = "IDR"
)

// Time layouts for parsing provider responses.
const (
	TimeLayoutRFC3339    = "2006-01-02T15:04:05Z07:00" // e.g. Garuda, AirAsia
	TimeLayoutISO8601    = "2006-01-02T15:04:05-0700"  // e.g. Batik Air
	TimeLayoutNoTimezone = "2006-01-02T15:04:05"       // e.g. Lion Air
)

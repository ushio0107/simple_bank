package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	return currency == USD || currency == EUR || currency == CAD
}

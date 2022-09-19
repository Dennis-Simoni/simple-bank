package util

const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case EUR, USD, GBP:
		return true
	}
	return false
}

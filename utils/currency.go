package utils

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// return true if IsSupportedCurrency
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, CAD, EUR:
		return true
	}

	return false
}

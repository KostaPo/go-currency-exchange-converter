package exchangerate

import "currency-exchange-converter/internal/currency"

type ExchangeRate struct {
	ID             int               `json:"id"`
	BaseCurrency   currency.Currency `json:"baseCurrency"`
	TargetCurrency currency.Currency `json:"targetCurrency"`
	Rate           float64           `json:"rate"`
}

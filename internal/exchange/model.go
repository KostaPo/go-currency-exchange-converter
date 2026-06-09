package exchange

import "currency-exchange-converter/internal/currency"

type Exchange struct {
	BaseCurrency    currency.Currency `json:"baseCurrency"`
	TargetCurrency  currency.Currency `json:"targetCurrency"`
	Rate            float64           `json:"rate"`
	Amount          float64           `json:"amount"`
	ConvertedAmount float64           `json:"convertedAmount"`
}

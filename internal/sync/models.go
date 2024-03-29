package sync

type LatestExchangeRate struct {
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate"`
}

type LatestExchangeRates []LatestExchangeRate

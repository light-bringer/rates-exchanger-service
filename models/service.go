package models

type LatestExchangeRate struct {
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate"`
}

type RateStatistic struct {
	Currency string  `json:"currency"`
	MinRate  float64 `json:"min_rate"`
	MaxRate  float64 `json:"max_rate"`
	AvgRate  float64 `json:"avg_rate"`
}

type (
	LatestExchangeRates []LatestExchangeRate
	RateStatistics      []RateStatistic
)

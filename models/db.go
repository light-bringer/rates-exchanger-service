package models

import "time"

type ExchangeRate struct {
	Currency string    `json:"currency"`
	Rate     float64   `json:"rate"`
	Time     time.Time `json:"time"`
}

type ExchangeRates []ExchangeRate

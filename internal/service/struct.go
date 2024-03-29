package service

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type RatesService struct {
	db        *pgxpool.Pool
	schema    string
	tableName string
}

// NewRatesService returns a new instance of RatesService.
func NewRatesService(db *pgxpool.Pool, schema string) *RatesService {
	return &RatesService{
		db:        db,
		schema:    schema,
		tableName: schema + ".exchange_rates",
	}
}

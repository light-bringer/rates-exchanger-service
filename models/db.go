package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type ExchangeRate struct {
	Currency string    `json:"currency"`
	Rate     float64   `json:"rate"`
	Time     time.Time `json:"time"`
}

type ExchangeRates []ExchangeRate

// PostgresConfig contains the configuration settings for PostgreSQL.
type PostgresConfig struct {
	Host           string `json:"host"            validate:"required"`
	Username       string `json:"username"        validate:"required"`
	Password       string `json:"password"        validate:"required"`
	Port           uint16 `json:"port"            validate:"required"`
	Database       string `json:"database"        validate:"required"`
	SSLMode        string `json:"ssl_mode"        validate:"required"`
	MaxConnections uint32 `json:"max_connections" validate:"required"`
	MinConnections uint32 `json:"min_connections" validate:"required"`
	SchemaName     string `json:"schema_name"     validate:"required"`
}

// Validate checks if the PostgresConfig is valid.
func (pc *PostgresConfig) Validate() bool {
	validate := validator.New()
	err := validate.Struct(pc)
	return err == nil
}

// PostgresConfigParams contains the configuration settings for PostgreSQL.
type PostgresConfigParams struct {
	Host           string `json:"host"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Port           uint16 `json:"port"`
	Database       string `json:"database"`
	SSLMode        string `json:"ssl_mode"`
	MaxConnections uint32 `json:"max_connections"`
	MinConnections uint32 `json:"min_connections"`
	SchemaName     string `json:"schema_name"`
}

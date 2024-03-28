package db

import (
	"github.com/go-playground/validator/v10"
)

const (
	SchemaName = "rate_api"
	TableName  = "exchange_rates"
)

// PostgresConfig contains the configuration settings for PostgreSQL.
type PostgresConfig struct {
	Host           string `json:"host"            validate:"required"`
	Username       string `json:"username"        validate:"required"`
	Password       string `json:"password"        validate:"required"`
	Port           uint16 `json:"port"            validate:"required"`
	Database       string `json:"database"        validate:"required"`
	SSLMode        string `json:"ssl_mode"        validate:"required"`
	MaxConnections int    `json:"max_connections" validate:"required"`
	MinConnections int    `json:"min_connections" validate:"required"`
	SchemaName     string `json:"schema_name"     validate:"required"`
}

// Validate checks if the PostgresConfig is valid.
func (pc *PostgresConfig) Validate() bool {
	validate := validator.New()
	err := validate.Struct(pc)
	return err == nil
}

type PostgresConfigParams struct {
	Host           string `json:"host"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Port           uint16 `json:"port"`
	Database       string `json:"database"`
	SSLMode        string `json:"ssl_mode"`
	MaxConnections int    `json:"max_connections"`
	MinConnections int    `json:"min_connections"`
	SchemaName     string `json:"schema_name"`
}

// NewPostgresConfig creates a new instance of PostgresConfig with the provided parameters.
// If the SSLMode is not specified, it defaults to "disable".
// If the SchemaName is not specified, it defaults to the value of the constant SchemaName.
// Returns a pointer to the created PostgresConfig instance if the configuration is valid, otherwise returns nil.
func NewPostgresConfig(params PostgresConfigParams) *PostgresConfig {
	if params.SSLMode == "" {
		params.SSLMode = "disable"
	}

	if params.SchemaName == "" {
		params.SchemaName = SchemaName
	}

	cfg := &PostgresConfig{
		Host:           params.Host,
		Username:       params.Username,
		Password:       params.Password,
		Port:           params.Port,
		Database:       params.Database,
		SSLMode:        params.SSLMode,
		MaxConnections: params.MaxConnections,
		MinConnections: params.MinConnections,
		SchemaName:     params.SchemaName,
	}
	if !cfg.Validate() {
		return nil
	}
	return cfg
}

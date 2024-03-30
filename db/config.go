package db

import "github.com/light-bringer/rates-exchanger-service/models"

const (
	SchemaName = "rate_api"
	TableName  = "exchange_rates"
)

// NewPostgresConfig creates a new instance of PostgresConfig with the provided parameters.
// If the SSLMode is not specified, it defaults to "disable".
// If the SchemaName is not specified, it defaults to the value of the constant SchemaName.
// Returns a pointer to the created PostgresConfig instance if the configuration is valid, otherwise returns nil.
func NewPostgresConfig(params models.PostgresConfigParams) *models.PostgresConfig {
	if params.SSLMode == "" {
		params.SSLMode = "disable"
	}

	if params.SchemaName == "" {
		params.SchemaName = SchemaName
	}

	cfg := &models.PostgresConfig{
		Host:           params.Host,
		Username:       params.Username,
		Password:       params.Password,
		Port:           params.Port,
		Database:       params.Database,
		SSLMode:        params.SSLMode,
		MaxConnections: uint32(params.MaxConnections),
		MinConnections: uint32(params.MinConnections),
		SchemaName:     params.SchemaName,
	}
	if !cfg.Validate() {
		return nil
	}
	return cfg
}

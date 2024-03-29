package db

import (
	"testing"

	"github.com/light-bringer/rates-exchanger-service/models"
	"github.com/stretchr/testify/assert"
)

func TestPostgresConfigValidate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		assert.True(t, (&models.PostgresConfig{
			Host:           "localhost",
			Username:       "user",
			Password:       "password",
			Port:           5432,
			Database:       "db",
			SSLMode:        "disable",
			MaxConnections: 10,
			MinConnections: 1,
			SchemaName:     "rate_api",
		}).Validate())
	})

	t.Run("invalid config", func(t *testing.T) {
		// Create a new instance of PostgresConfig
		pc := &models.PostgresConfig{}

		// Call the Validate method
		valid := pc.Validate()

		// Check if the validation is successful
		assert.False(t, valid)
	})
}

func TestNewPostgresConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		params := models.PostgresConfigParams{
			Host:           "localhost",
			Username:       "user",
			Password:       "password",
			Port:           5432,
			Database:       "db",
			SSLMode:        "disable",
			MaxConnections: 10,
			MinConnections: 1,
			SchemaName:     "rate_api",
		}

		cfg := NewPostgresConfig(params)

		assert.NotNil(t, cfg)
	})

	t.Run("invalid config", func(t *testing.T) {
		params := models.PostgresConfigParams{}

		cfg := NewPostgresConfig(params)

		assert.Nil(t, cfg)
	})

	t.Run("missing config", func(t *testing.T) {
		params := models.PostgresConfigParams{
			Host:           "localhost",
			Username:       "user",
			Password:       "password",
			Port:           5432,
			Database:       "db",
			SSLMode:        "disable",
			MaxConnections: 10,
			MinConnections: 1,
		}

		cfg := NewPostgresConfig(params)

		assert.NotNil(t, cfg)
	})
}

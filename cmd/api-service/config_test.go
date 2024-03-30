package main

import (
	"os"
	"testing"
	"time"

	"github.com/light-bringer/rates-exchanger-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	t.Run("valid config file", func(t *testing.T) {
		// Create a temporary config file for testing
		file, err := os.CreateTemp("", "config.yaml")
		require.NoError(t, err)
		defer os.Remove(file.Name())

		// Write valid config data to the file
		_, err = file.WriteString(`
database:
  host: localhost
  port: 8080
  db: db
`)
		require.NoError(t, err)

		// Call the readConfig function with the temporary file path
		config, err := readConfig(file.Name())

		// Check if the function returns the expected config and no error
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "localhost", config.Database.Host)
		assert.Equal(t, uint16(8080), config.Database.Port)
		assert.Equal(t, "db", config.Database.Name)
	})

	t.Run("non existent config file", func(t *testing.T) {
		// Call the readConfig function with a non-existent file path
		config, err := readConfig("nonexistent.yaml")

		// Check if the function returns an error and nil config
		require.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("invalid config file", func(t *testing.T) {
		// Create a temporary config file for testing
		file, err := os.CreateTemp("", "config.yaml")
		require.NoError(t, err)
		defer os.Remove(file.Name())

		// Write invalid config data to the file
		_, err = file.WriteString(`
database:
	host: localhost
	port: 8080
	db: db
	sslmode: disable
`)
		require.NoError(t, err)

		// Call the readConfig function with the temporary file path
		config, err := readConfig(file.Name())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal config")
		assert.Nil(t, config)

	})
}

func TestSetDefaults(t *testing.T) {
	config := &models.StartupConfig{
		Database: struct {
			Host           string         "yaml:\"host\""
			Port           uint16         "yaml:\"port\""
			User           string         "yaml:\"user\""
			Pass           string         "yaml:\"pass\""
			Name           string         "yaml:\"db\""
			Schema         string         "yaml:\"schema\""
			SSLMode        models.SSLMode "yaml:\"sslmode\""
			MinConnections uint32         "yaml:\"min_connections\""
			MaxConnections uint32         "yaml:\"max_connections\""
		}{
			Host:           "",
			Port:           0,
			User:           "",
			Pass:           "",
			Name:           "",
			Schema:         "",
			SSLMode:        "",
			MinConnections: 0,
			MaxConnections: 0,
		},
		CronJobs: struct {
			Rates struct {
				Enabled        bool          "yaml:\"enabled\""
				UpdateInterval time.Duration "yaml:\"interval\""
				SyncURL        string        "yaml:\"sync_url\""
			} "yaml:\"rates\""
			Cleanup struct {
				Enabled          bool          "yaml:\"enabled\""
				DeletionInterval time.Duration "yaml:\"interval\""
				MaxAge           int           "yaml:\"max_age\""
			} "yaml:\"cleanup\""
		}{
			Rates: struct {
				Enabled        bool          "yaml:\"enabled\""
				UpdateInterval time.Duration "yaml:\"interval\""
				SyncURL        string        "yaml:\"sync_url\""
			}{
				Enabled:        false,
				UpdateInterval: 0,
				SyncURL:        "",
			},
			Cleanup: struct {
				Enabled          bool          "yaml:\"enabled\""
				DeletionInterval time.Duration "yaml:\"interval\""
				MaxAge           int           "yaml:\"max_age\""
			}{
				Enabled:          false,
				DeletionInterval: 0,
				MaxAge:           0,
			},
		},
		HTTP: struct {
			Port int "yaml:\"port\""
		}{
			Port: 0,
		},
	}

	setDefaults(config)

	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, uint16(5432), config.Database.Port)
	assert.Equal(t, "postgres", config.Database.User)
	assert.Equal(t, "password", config.Database.Pass)
	assert.Equal(t, "rates", config.Database.Name)
	assert.Equal(t, "public", config.Database.Schema)
	assert.Equal(t, models.Disable, config.Database.SSLMode)
	assert.Equal(t, uint32(1), config.Database.MinConnections)
	assert.Equal(t, uint32(10), config.Database.MaxConnections)
	assert.Equal(t, syncURL, config.CronJobs.Rates.SyncURL)
	assert.Equal(t, syncInterval, config.CronJobs.Rates.UpdateInterval)
	assert.Equal(t, deleteInterval, config.CronJobs.Cleanup.DeletionInterval)
	assert.Equal(t, deletionDays, config.CronJobs.Cleanup.MaxAge)
	assert.Equal(t, 8080, config.HTTP.Port)
}

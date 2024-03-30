package main

import (
	"flag"
	"os"
	"time"

	"github.com/light-bringer/rates-exchanger-service/models"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	syncURL        = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	syncInterval   = 15 * time.Second
	deleteInterval = 1 * time.Minute
	ServerTimeout  = 15 * time.Second
	deletionDays   = 30
	contextTimeout = 60 * time.Second
)

// ReadConfig reads the configuration file from the given path and returns the StartupConfig.
func readConfig(path string) (*models.StartupConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	data = []byte(os.ExpandEnv(string(data)))

	// Unmarshal byte array into struct
	var config models.StartupConfig
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	return &config, nil
}

// setDefaults sets the default values for the configuration.
// If the configuration file does not have a value for a field, the default value is set.
func setDefaults(config *models.StartupConfig) {
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.User == "" {
		config.Database.User = "postgres"
	}

	if config.Database.Pass == "" {
		config.Database.Pass = "password"
	}

	if config.Database.Name == "" {
		config.Database.Name = "rates"
	}

	if config.Database.Schema == "" {
		config.Database.Schema = "public"
	}

	if config.Database.SSLMode == "" {
		config.Database.SSLMode = models.Disable
	}

	if config.Database.MinConnections == 0 {
		config.Database.MinConnections = 1
	}

	if config.Database.MaxConnections == 0 {
		config.Database.MaxConnections = 10
	}

	if config.CronJobs.Rates.SyncURL == "" {
		config.CronJobs.Rates.SyncURL = syncURL
	}
	if config.CronJobs.Rates.UpdateInterval == 0 {
		config.CronJobs.Rates.UpdateInterval = syncInterval
	}
	if config.CronJobs.Cleanup.DeletionInterval == 0 {
		config.CronJobs.Cleanup.DeletionInterval = deleteInterval
	}
	if config.CronJobs.Cleanup.MaxAge == 0 {
		config.CronJobs.Cleanup.MaxAge = deletionDays
	}
	if config.HTTP.Port == 0 {
		config.HTTP.Port = 8080
	}
}

// ParseFlags parses command-line flags into an AppConfig struct and returns it
func parseFlags() (string, error) {
	// Define flags

	dbHost := flag.String("config-file", "config.yaml", "path to the configuration file")
	// Parse the flags
	flag.Parse()

	if *dbHost == "" {
		return "", errors.New("dbhost flag is empty")
	}
	return *dbHost, nil

}

package models

import "time"

type StartupConfig struct {
	Database struct {
		Host           string  `yaml:"host"`
		Port           uint16  `yaml:"port"`
		User           string  `yaml:"user"`
		Pass           string  `yaml:"pass"`
		Name           string  `yaml:"db"`
		Schema         string  `yaml:"schema"`
		SSLMode        SSLMode `yaml:"sslmode"`
		MinConnections uint32  `yaml:"min_connections"`
		MaxConnections uint32  `yaml:"max_connections"`
	} `yaml:"database"`

	CronJobs struct {
		Rates struct {
			Enabled        bool          `yaml:"enabled"`
			UpdateInterval time.Duration `yaml:"interval"`
			SyncURL        string        `yaml:"sync_url"`
		} `yaml:"rates"`
		Cleanup struct {
			Enabled          bool          `yaml:"enabled"`
			DeletionInterval time.Duration `yaml:"interval"`
			MaxAge           int           `yaml:"max_age"`
		} `yaml:"cleanup"`
	} `yaml:"cronjobs"`

	HTTP struct {
		Port int `yaml:"port"`
	} `yaml:"http"`
}

type SSLMode string

const (
	Disable    SSLMode = "disable"
	Require    SSLMode = "require"
	VerifyFull SSLMode = "verify-full"
	VerifyCA   SSLMode = "verify-ca"
)

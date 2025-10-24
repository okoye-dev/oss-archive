package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Logging  LoggingConfig  `yaml:"logging"`
	S3       S3Config       `yaml:"s3"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port            int `yaml:"port"`
	ReadTimeout     int `yaml:"read_timeout"`     // in seconds
	WriteTimeout    int `yaml:"write_timeout"`    // in seconds
	IdleTimeout     int `yaml:"idle_timeout"`     // in seconds
	ShutdownTimeout int `yaml:"shutdown_timeout"` // in seconds
}

// LoggingConfig holds logging settings
type LoggingConfig struct {
	Level string `yaml:"level"` // debug, info, warn, error
	Mode  string `yaml:"mode"`  // gin mode: debug, release, test
}

// S3Config holds S3-compatible storage settings
type S3Config struct {
	Endpoint        string `yaml:"endpoint"`
	Region          string `yaml:"region"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	UseSSL          bool   `yaml:"use_ssl"`
	BucketName      string `yaml:"bucket_name"`
	ForcePathStyle  bool   `yaml:"force_path_style"`
}

// LoadConfig reads and parses the configuration file
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// GetDatabaseConnectionString returns a formatted database connection string
func (c *Config) GetDatabaseConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

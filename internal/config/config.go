package config

import (
	"fmt"
	"os"
	"strconv"

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
	// Try to load from environment variables first (for Railway/production)
	if config := loadFromEnv(); config != nil {
		return config, nil
	}

	// Fallback to config file (for local development)
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

func loadFromEnv() *Config {
	// Always try to load from env vars if any are set
	// No longer require database vars to be set

	// Log environment variables for debugging
	fmt.Printf("🔧 Loading configuration from environment:\n")
	fmt.Printf("   S3_ENDPOINT: %s\n", os.Getenv("S3_ENDPOINT"))
	fmt.Printf("   S3_REGION: %s\n", os.Getenv("S3_REGION"))
	fmt.Printf("   S3_ACCESS_KEY_ID: %s\n", maskString(os.Getenv("S3_ACCESS_KEY_ID")))
	fmt.Printf("   S3_SECRET_ACCESS_KEY: %s\n", maskString(os.Getenv("S3_SECRET_ACCESS_KEY")))
	fmt.Printf("   S3_USE_SSL: %s\n", os.Getenv("S3_USE_SSL"))
	fmt.Printf("   S3_BUCKET_NAME: %s\n", os.Getenv("S3_BUCKET_NAME"))
	fmt.Printf("   S3_FORCE_PATH_STYLE: %s\n", os.Getenv("S3_FORCE_PATH_STYLE"))
	fmt.Printf("   PORT: %s\n", os.Getenv("PORT"))

	config := &Config{
		// Database: DatabaseConfig{
		// 	Host:     getEnv("DB_HOST", "localhost"),
		// 	Port:     getEnvInt("DB_PORT", 5432),
		// 	User:     getEnv("DB_USER", "postgres"),
		// 	Password: getEnv("DB_PASSWORD", ""),
		// 	DBName:   getEnv("DB_NAME", "oss-archive"),
		// 	SSLMode:  getEnv("DB_SSLMODE", "require"),
		// },
		Server: ServerConfig{
			Port:            getEnvInt("PORT", 6060),
			ReadTimeout:     getEnvInt("READ_TIMEOUT", 30),
			WriteTimeout:    getEnvInt("WRITE_TIMEOUT", 30),
			IdleTimeout:     getEnvInt("IDLE_TIMEOUT", 120),
			ShutdownTimeout: getEnvInt("SHUTDOWN_TIMEOUT", 5),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			Mode:  getEnv("GIN_MODE", "release"),
		},
		S3: S3Config{
			Endpoint:        getEnv("S3_ENDPOINT", ""),
			Region:          getEnv("S3_REGION", "us-east-1"),
			AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),
			UseSSL:          getEnvBool("S3_USE_SSL", true),
			BucketName:      getEnv("S3_BUCKET_NAME", "oss-archive"),
			ForcePathStyle:  getEnvBool("S3_FORCE_PATH_STYLE", false),
		},
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// maskString masks sensitive strings for logging
func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
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

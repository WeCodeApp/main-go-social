package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Services ServicesConfig
	Logging  LoggingConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            int
	Username        string
	Password        string
	Name            string
	Charset         string
	ParseTime       bool
	Loc             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// ServicesConfig holds URLs for other microservices
type ServicesConfig struct {
	UsersServiceURL   string
	FriendsServiceURL string
	GroupsServiceURL  string
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level  string
	Format string
	Output string
	File   string
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s&charset=%s",
		c.Username, c.Password, c.Host, c.Port, c.Name, c.Loc, c.Charset)
}
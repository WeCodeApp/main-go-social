package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the gateway API
type Config struct {
	Environment string `mapstructure:"environment"`
	Port        string `mapstructure:"port"`

	// gRPC client configurations
	UsersServiceURL   string `mapstructure:"users_service_url"`
	PostsServiceURL   string `mapstructure:"posts_service_url"`
	FriendsServiceURL string `mapstructure:"friends_service_url"`
	GroupsServiceURL  string `mapstructure:"groups_service_url"`

	// App URL
	AppURL string `mapstructure:"app_url"`

	// Auth configurations
	JWTSecret string `mapstructure:"jwt_secret"`

	// OAuth configurations
	OAuth struct {
		Google struct {
			ClientID     string `mapstructure:"client_id"`
			ClientSecret string `mapstructure:"client_secret"`
			RedirectURL  string `mapstructure:"redirect_url"`
		} `mapstructure:"google"`
		Microsoft struct {
			ClientID     string `mapstructure:"client_id"`
			ClientSecret string `mapstructure:"client_secret"`
			RedirectURL  string `mapstructure:"redirect_url"`
		} `mapstructure:"microsoft"`
	} `mapstructure:"oauth"`

	// Logging configurations
	LogLevel string `mapstructure:"log_level"`
}

// LoadConfig loads configuration from environment variables and config files
func LoadConfig() (*Config, error) {
	// Set default values
	viper.SetDefault("environment", "development")
	viper.SetDefault("port", "8000")
	viper.SetDefault("users_service_url", "localhost:50051")
	viper.SetDefault("posts_service_url", "localhost:50052")
	viper.SetDefault("friends_service_url", "localhost:50053")
	viper.SetDefault("groups_service_url", "localhost:50054")
	viper.SetDefault("jwt_secret", "your-secret-key")
	viper.SetDefault("log_level", "info")

	// OAuth default values
	viper.SetDefault("oauth.google.client_id", "your-google-client-id")
	viper.SetDefault("oauth.google.client_secret", "your-google-client-secret")
	viper.SetDefault("oauth.google.redirect_url", "http://localhost:8000/api/v1/auth/google/callback")
	viper.SetDefault("oauth.microsoft.client_id", "your-microsoft-client-id")
	viper.SetDefault("oauth.microsoft.client_secret", "your-microsoft-client-secret")
	viper.SetDefault("oauth.microsoft.redirect_url", "http://localhost:8000/api/v1/auth/microsoft/callback")

	// Set config file name and paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.social-media")

	// Read environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(viper.ConfigFileUsed())
	if configDir == "." {
		configDir = "./config"
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return nil, err
			}
		}
	}

	// Write default config file if it doesn't exist
	if viper.ConfigFileUsed() == "" {
		defaultConfig := map[string]interface{}{
			"environment":         config.Environment,
			"port":                config.Port,
			"users_service_url":   config.UsersServiceURL,
			"posts_service_url":   "127.0.0.1:50052",
			"friends_service_url": config.FriendsServiceURL,
			"groups_service_url":  config.GroupsServiceURL,
			"jwt_secret":          config.JWTSecret,
			"log_level":           config.LogLevel,
			"oauth": map[string]interface{}{
				"google": map[string]interface{}{
					"client_id":     "your-google-client-id",
					"client_secret": "your-google-client-secret",
					"redirect_url":  "http://localhost:8000/api/v1/auth/google/callback",
				},
				"microsoft": map[string]interface{}{
					"client_id":     "your-microsoft-client-id",
					"client_secret": "your-microsoft-client-secret",
					"redirect_url":  "http://localhost:8000/api/v1/auth/microsoft/callback",
				},
			},
		}

		configFile := filepath.Join(configDir, "config.yaml")
		jsonData, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return nil, err
		}

		if err := os.WriteFile(configFile, jsonData, 0644); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

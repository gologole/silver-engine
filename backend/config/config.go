package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (Config, error) {
	config := Config{
		DBHost:     os.Getenv("DBHost"),
		DBPort:     os.Getenv("DBPort"),
		DBUser:     os.Getenv("DBUser"),
		DBPassword: os.Getenv("DBPassword"),
		DBName:     os.Getenv("DBName"),
		AppPort:    os.Getenv("AppPort"),
	}

	// Validate required environment variables
	missingVars := []string{}
	if config.DBHost == "" {
		missingVars = append(missingVars, "DBHost")
	}
	if config.DBPort == "" {
		missingVars = append(missingVars, "DBPort")
	}
	if config.DBUser == "" {
		missingVars = append(missingVars, "DBUser")
	}
	if config.DBPassword == "" {
		missingVars = append(missingVars, "DBPassword")
	}
	if config.DBName == "" {
		missingVars = append(missingVars, "DBName")
	}
	if config.AppPort == "" {
		missingVars = append(missingVars, "AppPort")
	}

	if len(missingVars) > 0 {
		return config, fmt.Errorf("missing required environment variables: %v", missingVars)
	}

	return config, nil
}

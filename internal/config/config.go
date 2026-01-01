package config

import "os"

type Config struct {
	DatabaseURL string
	Port        string // Add Port field
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/smartq_db?sslmode=disable"),
		Port:        getEnv("PORT", "8080"), // Initialize Port
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

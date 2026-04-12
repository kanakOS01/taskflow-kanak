package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	return &Config{
		Port:        getEnv("PORT", "8000"),
		DatabaseURL: getRequiredEnv("DATABASE_URL"),
		JWTSecret:   getRequiredEnv("JWT_SECRET"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getRequiredEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		panic(fmt.Sprintf("missing required environment variable: %s", key))
	}
	return value
}

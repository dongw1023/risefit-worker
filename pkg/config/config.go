package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port                string
	InternalAPIKey      string
	EmailProviderAPIKey string
	FromEmail           string
}

func Load() (*Config, error) {
	port := getEnv("PORT", "8080")
	internalAPIKey := os.Getenv("INTERNAL_API_KEY")
	emailProviderAPIKey := os.Getenv("EMAIL_PROVIDER_API_KEY")
	fromEmail := os.Getenv("FROM_EMAIL")

	if internalAPIKey == "" {
		return nil, fmt.Errorf("INTERNAL_API_KEY is required")
	}
	if emailProviderAPIKey == "" {
		return nil, fmt.Errorf("EMAIL_PROVIDER_API_KEY is required")
	}
	if fromEmail == "" {
		return nil, fmt.Errorf("FROM_EMAIL is required")
	}

	return &Config{
		Port:                port,
		InternalAPIKey:      internalAPIKey,
		EmailProviderAPIKey: emailProviderAPIKey,
		FromEmail:           fromEmail,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

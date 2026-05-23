package main

import (
	"os"
	"strings"
	"time"
)

type config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	JWTIssuer      string
	JWTTTL         time.Duration
	KafkaBrokers   []string
	KafkaUserTopic string
	KafkaEnabled   bool
}

func loadConfig() config {
	return config{
		Port:           env("PORT", "8080"),
		DatabaseURL:    env("DATABASE_URL", "postgres://users_app:users_pass@users-db:5432/users?sslmode=disable"),
		JWTSecret:      env("JWT_SECRET", "dev-user-service-secret-change-me"),
		JWTIssuer:      env("JWT_ISSUER", "patitomedi-user-service"),
		JWTTTL:         durationEnv("JWT_TTL", 24*time.Hour),
		KafkaBrokers:   splitEnv("KAFKA_BROKERS", "kafka:9092"),
		KafkaUserTopic: env("KAFKA_USER_TOPIC", "user-registered"),
		KafkaEnabled:   env("KAFKA_ENABLED", "true") == "true",
	}
}

func env(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func splitEnv(key, fallback string) []string {
	value := env(key, fallback)
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

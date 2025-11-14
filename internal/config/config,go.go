package config

import (
	"os"
)

type Config struct {
	HTTPPort string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
	DBSSL    string
}

func Load() *Config {
	cfg := &Config{
		HTTPPort: getEnv("HTTP_PORT", ":8080"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5432"),
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPass:   getEnv("DB_PASS", "postgres"),
		DBName:   getEnv("DB_NAME", "qna"),
		DBSSL:    getEnv("DB_SSLMODE", "disable"),
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

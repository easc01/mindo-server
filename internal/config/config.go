package config

import (
	"os"

	"github.com/ishantSikdar/mindo-server/pkg/logger"
	"github.com/joho/godotenv"
)

func init() {
	dotenvErr := godotenv.Load()
	if dotenvErr != nil {
		logger.Log.Errorf(
			"error loading .env file, continuing with system environment variables, %s",
			dotenvErr,
		)
	}
}

type Config struct {
	ApiBaseUrl         string
	DbHost             string
	DbPort             string
	DbUser             string
	DbPassword         string
	DbName             string
	JwtSecret          string
	GoogleAPIKey       string
	GoogleClientId     string
	GoogleClientSecret string
}

func GetConfig() *Config {
	return &Config{
		DbHost:             getEnv("DB_HOST", "localhost"),
		DbPort:             getEnv("DB_PORT", "5432"),
		DbUser:             getEnv("DB_USER", "postgres"),
		DbPassword:         getEnv("DB_PASSWORD", "__DB_PASSWORD__"),
		DbName:             getEnv("DB_NAME", "mindo"),
		JwtSecret:          getEnv("JWT_SECRET", "__JWT_SECRET__"),
		GoogleAPIKey:       getEnv("GOOGLE_API_KEY", "__GOOGLE_API_KEY__"),
		GoogleClientId:     getEnv("GOOGLE_CLIENT_ID", "__GOOGLE_CLIENT_ID__"),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", "__GOOGLE_CLIENT_SECRET__"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

package config

import (
	"os"

	"github.com/easc01/mindo-server/pkg/logger"
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

type Environment string

const (
	Development Environment = "dev"
	Production  Environment = "prod"
)

type Config struct {
	AppPort            string
	Env                Environment
	DbConnectionUri    string
	JwtSecret          string
	GoogleAPIKey       string
	GoogleClientId     string
	GoogleClientSecret string
	YoutubeAPIKey      string
}

func GetConfig() *Config {
	return &Config{
		DbConnectionUri:    getEnv("DB_URI", "postgres://postgres:6515@localhost:5432/mindo"),
		AppPort:            getEnv("APP_PORT", "8080"),
		Env:                Environment(getEnv("ENV", "dev")),
		JwtSecret:          getEnv("JWT_SECRET", "__JWT_SECRET__"),
		GoogleAPIKey:       getEnv("GOOGLE_API_KEY", "__GOOGLE_API_KEY__"),
		GoogleClientId:     getEnv("GOOGLE_CLIENT_ID", "__GOOGLE_CLIENT_ID__"),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", "__GOOGLE_CLIENT_SECRET__"),
		YoutubeAPIKey:      getEnv("YOUTUBE_API_KEY", "__YOUTUBE_API_KEY__"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

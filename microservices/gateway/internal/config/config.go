package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	CommandServiceURL string
	QueryServiceURL   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading configuration from environment variables")
	}

	return &Config{
		Port:              getEnv("PORT"),
		CommandServiceURL: getEnv("COMMAND_SERVICE_URL"),
		QueryServiceURL:   getEnv("QUERY_SERVICE_URL"),
	}
}

func getEnv(key string) string {
	val, ok := os.LookupEnv(key)

	if !ok || val == "" {
		log.Fatalf("missing required env: %s", key)
	}

	return val
}

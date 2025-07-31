package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	KafkaBrokers   string
	KafkaTopic     string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
}

// Load loads the configuration from environment variables
func Load() *Config {
	// load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		KafkaBrokers:   getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:     getEnv("KAFKA_TOPIC", "stock-data-raw"),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "username"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "password"),
	}
}

// getEnv retrieves an environment variable or returns a fallback value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

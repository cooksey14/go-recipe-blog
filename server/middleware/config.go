package middleware

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Function to load environment variables
func LoadEnvConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or unable to load it")
	}
}

// Function to get required environment variables
func GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalf("Environment variable %s is not set", key)
	return ""
}

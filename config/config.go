package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	JWTSecret       string
	OpenLibBaseURL  string
	OpenLibImageURL string
	OpenLibLimit    int
	SMTPHost        string
	SMTPPort        int
	SMTPUsername    string
	SMTPPassword    string
	EmailFrom       string
}

var Config *Configuration

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found. Falling back to system environment variables.")
	}

	Config = &Configuration{
		JWTSecret:       GetEnv("JWT_SECRET", "jwt_secret"),
		OpenLibBaseURL:  GetEnv("OPEN_LIB_BASE_URL", "https://openlibrary.org/subjects/%s.json?limit=%d&offset=%d"),
		OpenLibImageURL: GetEnv("OPEN_LIB_IMAGE_URL", "https://covers.openlibrary.org/b/id/%d-L.jpg"),
		OpenLibLimit:    GetEnvAsInt("OPEN_LIB_LIMIT", 12),
		SMTPHost:        GetEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:        GetEnvAsInt("SMTP_PORT", 587),
		SMTPUsername:    GetEnv("SMTP_USERNAME", ""),
		SMTPPassword:    GetEnv("SMTP_PASSWORD", ""),
		EmailFrom:       GetEnv("EMAIL_FROM", "no-reply@library.com"),
	}
}

func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == "" {
			log.Fatalf("Environment variable %s is required but not set", key)
		}
		return defaultValue
	}
	return value
}

func GetEnvAsInt(key string, defaultValue int) int {
	value := GetEnv(key, fmt.Sprintf("%d", defaultValue))
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid value for %s, expected integer, got %s", key, value)
	}
	return intValue
}

package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerPort string

	// PG
	JWTSecret string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// mail
	SMTPHost     string
    SMTPPort     int
    SMTPUser     string
    SMTPPassword string
    SMTPFromEmail string
    SMTPFromName  string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		// Server
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "apartment_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		//JWT
		JWTSecret:  getEnv("JWT_SECRET", "fairytail"),

		// Mail
		SMTPHost:      getEnv("SMTP_HOST", "localhost"),
		SMTPPort:      getEnvAsInt("SMTP_PORT", 1025),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
		SMTPFromEmail: getEnv("SMTP_FROM_EMAIL", "no-reply@rental.com"),
		SMTPFromName:  getEnv("SMTP_FROM_NAME", "Rental Service"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func (c *Config) GetDBConnectionString() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

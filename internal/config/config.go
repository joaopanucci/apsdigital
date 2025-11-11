package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	Upload   UploadConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret                 string
	Expiration            string
	RefreshTokenExpiration string
}

type ServerConfig struct {
	Port string
	Env  string
}

type UploadConfig struct {
	Path string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "apsdigital"),
		},
		JWT: JWTConfig{
			Secret:                 getEnv("JWT_SECRET", "your-super-secret-jwt-key-here"),
			Expiration:            getEnv("JWT_EXPIRATION", "24h"),
			RefreshTokenExpiration: getEnv("REFRESH_TOKEN_EXPIRATION", "168h"),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Upload: UploadConfig{
			Path: getEnv("UPLOAD_PATH", "./uploads"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
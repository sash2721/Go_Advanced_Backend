package configs

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port      string
	Env       string
	SecretKey string
}

func GetServerConfig() *ServerConfig {
	err := godotenv.Load()

	if err != nil {
		slog.Warn("Error loading .env file, using system environment variables")
	}

	serverConfig := &ServerConfig{
		Port:      os.Getenv("PORT"),
		Env:       os.Getenv("ENV"),
		SecretKey: os.Getenv("SECRET_KEY"),
	}

	return serverConfig
}

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppPort string

	DBUser string
	DBPass string
	DBName string
	DBHost string
	DBPort string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		AppName: os.Getenv("APP_NAME"),
		AppPort: os.Getenv("APP_PORT"),

		DBUser: os.Getenv("POSTGRES_USER"),
		DBPass: os.Getenv("POSTGRES_PASSWORD"),
		DBName: os.Getenv("POSTGRES_DB"),
		DBHost: os.Getenv("POSTGRES_HOST"),
		DBPort: os.Getenv("POSTGRES_PORT"),
	}

	if cfg.AppPort == "" {
		log.Fatal("APP_PORT not set")
	}

	return cfg
}

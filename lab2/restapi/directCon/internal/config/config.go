package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

var AppConfig *Config

func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		log.Fatalf("Failed on loading .env: %v", err)
	}

	AppConfig = &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME_EASY"),
		DBPort:     os.Getenv("DB_PORT"),
	}
}

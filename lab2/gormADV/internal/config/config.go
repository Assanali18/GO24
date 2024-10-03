package config

import (
	"github.com/go-playground/validator/v10"
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
var Validate = validator.New()

func init() {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatalf("Failed on loading .env: %v", err)
	}

	AppConfig = &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME_GORM"),
		DBPort:     os.Getenv("DB_PORT"),
	}
}

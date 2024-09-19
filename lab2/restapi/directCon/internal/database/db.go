package database

import (
	"database/sql"
	"directCon/internal/config"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}

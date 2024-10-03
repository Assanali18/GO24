package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID   uint `gorm:"primary_key"`
	Name string
	Age  int
}

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME_EASY")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	fmt.Println("Подключение к базе данных успешно")

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
	fmt.Println("Миграция успешно завершена")

	insertData(db)

	queryData(db)
}

func insertData(db *gorm.DB) {
	users := []User{
		{Name: "Asan", Age: 19},
		{Name: "Alisher", Age: 20},
		{Name: "Aibek", Age: 21},
	}
	result := db.Create(&users)
	if result.Error != nil {
		fmt.Println("Ошибка вставки данных: ", result.Error)
	} else {
		fmt.Println("Данные успешно вставлены")
	}
}

func queryData(db *gorm.DB) {
	var users []User

	result := db.Find(&users)
	if result.Error != nil {
		fmt.Println("Ошибка запроса данных: ", result.Error)
	} else {
		fmt.Printf("Found %d users\n", result.RowsAffected)
		for _, user := range users {
			fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
		}
	}
}

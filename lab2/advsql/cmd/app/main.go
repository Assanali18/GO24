package main

import (
	"advsql/internal/app"
	"advsql/internal/database"
	"fmt"
)

func main() {
	fmt.Println("Сервер запущен на порту 8080")
	db, err := database.ConnectDB()
	if err != nil {
		return
	} else {
		fmt.Println("Database connected")
	}
	database.DB = db
	app.Run()
}

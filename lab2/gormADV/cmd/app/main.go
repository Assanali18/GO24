package main

import (
	"fmt"
	"gormADV/internal/app"
	"gormADV/internal/database"
)

func main() {
	fmt.Println("Сервер запущен на порту 8080")
	database.ConnectDB()
	app.Run()
}

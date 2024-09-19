package main

import (
	"fmt"
	"gorm/internal/app"
	"gorm/internal/database"
)

func main() {
	fmt.Println("Сервер запущен на порту 8080")
	database.ConnectDB()
	app.Run()
}

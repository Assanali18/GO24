package main

import (
	"fmt"
	"gormADV/internal/app"
	"gormADV/internal/database"
)

// @title           GO REST API WITH GORM
// @version         1.0
// @description     This is a sample API server.
// @termsOfService  http://example.com/terms/

// @contact.name    API Support
// @contact.url     http://www.example.com/support
// @contact.email   support@example.com

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

func main() {
	fmt.Println("Сервер запущен на порту 8080")
	database.ConnectDB()
	app.Run()
}

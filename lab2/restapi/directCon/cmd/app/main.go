package main

import (
	"directCon/internal/app"
	"fmt"
)

func main() {
	fmt.Println("Сервер запущен на порту 8080")
	app.Run()
}

package app

import (
	"gormADV/internal/database"
	"gormADV/internal/models"
	"gormADV/internal/transport"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "gormADV/docs"
)

func Run() {
	err := database.DB.AutoMigrate(&models.User{}, &models.Profile{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
	log.Println("Auto migration completed.")

	r := mux.NewRouter()

	transport.RegisterRoutes(r)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}

package app

import (
	"advsql/internal/services"
	"advsql/internal/transport"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"

	_ "advsql/docs"
	"github.com/gorilla/mux"
)

func Run() {
	err := services.СreateUsersTable()
	if err != nil {
		return
	}
	r := mux.NewRouter()

	transport.RegisterRoutes(r)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}

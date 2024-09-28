package app

import (
	"advsql/internal/services"
	"advsql/internal/transport"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Run() {
	err := services.СreateUsersTable()
	if err != nil {
		return
	}
	r := mux.NewRouter()

	transport.RegisterRoutes(r)

	log.Fatal(http.ListenAndServe(":8080", r))
}

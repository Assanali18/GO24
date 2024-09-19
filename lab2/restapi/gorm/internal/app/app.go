package app

import (
	"gorm/internal/transport"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	transport.RegisterRoutes(r)

	log.Fatal(http.ListenAndServe(":8000", r))
}

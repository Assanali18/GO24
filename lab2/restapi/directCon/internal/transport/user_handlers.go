package transport

import (
	"directCon/internal/models"
	"directCon/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/user", CreateUser).Methods("POST")
	r.HandleFunc("/user/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE")
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := services.GetUsers()
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := services.CreateUser(user); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	user.ID = id

	userData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	log.Printf("Updating user: %s", userData)

	if err := services.UpdateUser(user); err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := services.DeleteUser(id); err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package transport

import (
	_ "advsql/docs"
	"advsql/internal/models"
	"advsql/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// RegisterRoutes registers all routes for the application.
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/users", CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", DeleteUser).Methods(http.MethodDelete)
}

// GetUsers	Get list of users
// @Summary Get list of users
// @Description Get a paginated list of users with optional filters for age and sorting in ascending or descending order.
// @Tags users
// @Accept  json
// @Produce  json
// @Param   min_age query int false "Minimum Age"
// @Param   max_age query int false "Maximum Age"
// @Param   page query int false "Page number"
// @Param   page_size query int false "Page size"
// @Param   sort query string false "Sort by name in ascending or descending order"
// @Success 200 {object} models.UserListResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	minAge, _ := strconv.Atoi(r.URL.Query().Get("min_age"))
	maxAge, _ := strconv.Atoi(r.URL.Query().Get("max_age"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	sort := r.URL.Query().Get("sort")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalCount, err := services.GetUsers(minAge, maxAge, page, pageSize, sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	response := struct {
		Users      []models.User `json:"users"`
		TotalItems int           `json:"total_items"`
		Page       int           `json:"page"`
		PageSize   int           `json:"page_size"`
		TotalPages int           `json:"total_pages"`
	}{
		Users:      users,
		TotalItems: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateUser creates a new users.
// @Summary     Create few users at once
// @Description Create a new users
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       user body     []models.User true "User to create"
// @Success     201  {string} string "Created"
// @Failure     400  {string} string "Invalid request payload"
// @Failure     500  {string} string "Internal server error"
// @Router      /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := services.CreateUser(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// UpdateUser updates an existing user.
// @Summary     Update user
// @Description Update user details by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id   path     int         true "User ID"
// @Param       user body     models.User true "Updated user"
// @Success     200  {object} models.User
// @Failure     400  {string} string "Invalid request payload"
// @Failure     404  {string} string "User not found"
// @Failure     500  {string} string "Internal server error"
// @Router      /users/{id} [put]
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(user)
}

// DeleteUser deletes a user.
// @Summary     Delete user
// @Description Delete user by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id  path     int     true "User ID"
// @Success     204 {string} string "No Content"
// @Failure     400 {string} string "Invalid user ID"
// @Failure     404 {string} string "User not found"
// @Failure     500 {string} string "Internal server error"
// @Router      /users/{id} [delete]
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

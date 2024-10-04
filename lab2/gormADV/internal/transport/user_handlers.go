package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gormADV/internal/config"
	"gormADV/internal/models"
	"gormADV/internal/services"
	"net/http"
	"strconv"
)

// RegisterRoutes registers all routes for the application.
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
}

// GetUsers @Summary Get list of users
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

	users, totalCount, err := services.GetUsersWithProfiles(minAge, maxAge, page, pageSize, sort)
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

// CreateUser creates a new user.
// @Summary     Create user
// @Description Create a new user with profile
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       user body     models.User true "User to create"
// @Success     201  {object} models.User
// @Failure     400  {string} string "Invalid request payload"
// @Failure     500  {string} string "Internal server error"
// @Router      /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := config.Validate.Struct(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := services.CreateUserWithProfile(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user.ID = uint(userID)

	if err := config.Validate.Struct(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := services.UpdateUserAndProfile(&user, user.Profile); err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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

	if err := services.DeleteUserWithProfile(uint(id)); err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package transport_test

import (
	"advsql/internal/database"
	"advsql/internal/models"
	"advsql/internal/transport"
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockDB sqlmock.Sqlmock

func setupMockDB(t *testing.T) {
	var db *sql.DB
	var err error
	db, mockDB, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось подключиться к sqlmock: %v", err)
	}
	database.DB = db
}

func TestGetUsers(t *testing.T) {
	setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(1, "John Doe", 25).
		AddRow(2, "Jane Doe", 30)

	mockDB.ExpectQuery("SELECT id, name, age FROM users").
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/users?page=1&page_size=10&min_age=18&max_age=30&sort=name_asc", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/users", transport.GetUsers)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неверный код статуса: получили %v, ожидали %v", status, http.StatusOK)
	}

	var users []models.User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Errorf("Ошибка при декодировании ответа: %v", err)
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

func TestCreateUser(t *testing.T) {
	setupMockDB(t)

	mockDB.ExpectBegin()
	mockDB.ExpectPrepare("INSERT INTO users").
		ExpectExec().
		WithArgs("John Doe", 25).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()

	users := []models.User{
		{Name: "John Doe", Age: 25},
	}
	payload, err := json.Marshal(users)
	if err != nil {
		t.Fatalf("Не удалось сериализовать пользователей: %v", err)
	}

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/users", transport.CreateUser)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Неверный код статуса: получили %v, ожидали %v", status, http.StatusCreated)
	}
}

func TestUpdateUser(t *testing.T) {
	setupMockDB(t)

	mockDB.ExpectExec("UPDATE users SET").
		WithArgs("Jane Doe", 30, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	user := models.User{Name: "Jane Doe", Age: 30}
	payload, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Не удалось сериализовать пользователя: %v", err)
	}

	req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", transport.UpdateUser)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Неверный код статуса: получили %v, ожидали %v", status, http.StatusOK)
	}

	var updatedUser models.User
	if err := json.NewDecoder(rr.Body).Decode(&updatedUser); err != nil {
		t.Errorf("Ошибка при декодировании ответа: %v", err)
	}

	if updatedUser.Name != "Jane Doe" || updatedUser.Age != 30 {
		t.Errorf("Неверные данные пользователя: получили %v", updatedUser)
	}
}

func TestDeleteUser(t *testing.T) {
	setupMockDB(t)

	mockDB.ExpectExec("DELETE FROM users WHERE").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("DELETE", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", transport.DeleteUser)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Неверный код статуса: получили %v, ожидали %v", status, http.StatusNoContent)
	}
}

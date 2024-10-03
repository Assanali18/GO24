package transport_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gormADV/internal/database"
	"gormADV/internal/models"
	"gormADV/internal/services"
	"testing"
)

var mockDB *gorm.DB
var mock sqlmock.Sqlmock

func setupMockDB(t *testing.T) {
	var db *sql.DB
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать mock базу данных: %v", err)
	}
	mockDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Не удалось подключиться к GORM с mock базой данных: %v", err)
	}

	database.DB = mockDB
}

func TestCreateUserWithProfile(t *testing.T) {
	setupMockDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" \("created_at","updated_at","deleted_at","name","age"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "John Doe", 25).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	user := &models.User{Name: "John Doe", Age: 25}
	err := services.CreateUserWithProfile(user)
	if err != nil {
		t.Errorf("Создание пользователя завершилось с ошибкой: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

func TestGetUsersWithProfiles(t *testing.T) {
	setupMockDB(t)

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "users" WHERE age >= \$1 AND age <= \$2 AND "users"\."deleted_at" IS NULL`).
		WithArgs(18, 30).
		WillReturnRows(countRows)

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(1, "John Doe", 25).
		AddRow(2, "Jane Doe", 30)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE age >= \$1 AND age <= \$2 AND "users"\."deleted_at" IS NULL ORDER BY name ASC LIMIT \$3`).
		WithArgs(18, 30, 10).
		WillReturnRows(rows)

	profileRows := sqlmock.NewRows([]string{"id", "user_id", "bio", "profile_picture_url"}).
		AddRow(1, 1, "Bio for John", "http://example.com/john.jpg").
		AddRow(2, 2, "Bio for Jane", "http://example.com/jane.jpg")

	mock.ExpectQuery(`SELECT \* FROM "profiles" WHERE "profiles"\."user_id" IN \(\$1,\$2\) AND "profiles"\."deleted_at" IS NULL`).
		WithArgs(1, 2).
		WillReturnRows(profileRows)

	users, totalCount, err := services.GetUsersWithProfiles(18, 30, 1, 10, "name_asc")
	if err != nil {
		t.Errorf("Получение списка пользователей завершилось с ошибкой: %v", err)
	}

	if totalCount != 2 {
		t.Errorf("Ожидалось общее количество пользователей 2, но получили %v", totalCount)
	}

	if len(users) != 2 {
		t.Errorf("Ожидалось 2 пользователя, но получили %v", len(users))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

func TestUpdateUserAndProfile(t *testing.T) {
	setupMockDB(t)

	user := &models.User{
		Model: gorm.Model{ID: 1},
		Name:  "Jane Doe",
		Age:   30,
	}
	profile := &models.Profile{
		UserID:            1,
		Bio:               "Some Profile Data",
		ProfilePictureURL: "http://example.com/profile.jpg",
	}

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "users" SET "age"=\$1,"name"=\$2,"updated_at"=\$3 WHERE id = \$4 AND "users"."deleted_at" IS NULL`).
		WithArgs(30, "Jane Doe", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`UPDATE "profiles" SET "bio"=\$1,"profile_picture_url"=\$2,"updated_at"=\$3 WHERE user_id = \$4 AND "profiles"."deleted_at" IS NULL`).
		WithArgs("Some Profile Data", "http://example.com/profile.jpg", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := services.UpdateUserAndProfile(user, profile)
	if err != nil {
		t.Errorf("Обновление пользователя завершилось с ошибкой: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

func TestDeleteUserWithProfile(t *testing.T) {
	setupMockDB(t)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "users" SET "deleted_at"=\$1 WHERE "users"."id" = \$2 AND "users"."deleted_at" IS NULL`).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`UPDATE "profiles" SET "deleted_at"=\$1 WHERE user_id = \$2 AND "profiles"."deleted_at" IS NULL`).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := services.DeleteUserWithProfile(1)
	if err != nil {
		t.Errorf("Удаление пользователя завершилось с ошибкой: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

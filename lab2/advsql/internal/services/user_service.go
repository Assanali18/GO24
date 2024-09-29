package services

import (
	"advsql/internal/database"
	"advsql/internal/models"
	"fmt"
	"log"
)

func СreateUsersTable() error {
	query := `
   CREATE TABLE IF NOT EXISTS users (
       id SERIAL PRIMARY KEY,
       name VARCHAR(255) UNIQUE NOT NULL,
       age INT NOT NULL
   );
   `
	if _, err := database.DB.Exec(query); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	log.Println("Users table created or already exists.")
	return nil
}

func CreateUser(users []models.User) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO users (name, age) VALUES ($1, $2)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, user := range users {
		if _, err := stmt.Exec(user.Name, user.Age); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func GetUsers(minAge, maxAge, page, pageSize int, sort string) ([]models.User, error) {
	offset := (page - 1) * pageSize

	query := `
    SELECT id, name, age FROM users
    WHERE 1=1
    `

	var params []interface{}
	index := 1

	if minAge > 0 {
		query += fmt.Sprintf(" AND age >= $%d", index)
		params = append(params, minAge)
		index++
	}
	if maxAge > 0 {
		query += fmt.Sprintf(" AND age <= $%d", index)
		params = append(params, maxAge)
		index++
	}

	switch sort {
	case "name_asc":

		query += " ORDER BY name ASC"
	case "name_desc":
		query += " ORDER BY name DESC"
	default:
		query += " ORDER BY id"
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", index, index+1)
	params = append(params, pageSize, offset)

	rows, err := database.DB.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}

func UpdateUser(user models.User) error {
	result, err := database.DB.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", user.Name, user.Age, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func DeleteUser(userID int) error {
	result, err := database.DB.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

package services

import (
	"directCon/internal/database"
	"directCon/internal/models"
	"errors"
)

func CreateUser(user models.User) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id
                              `
	err = db.QueryRow(query, user.Name, user.Age).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetUsers() ([]models.User, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
SELECT id, name, age FROM users
`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func UpdateUser(user models.User) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
UPDATE users SET name=$1, age=$2 WHERE id=$3
`
	result, err := db.Exec(query, user.Name, user.Age, user.ID)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func DeleteUser(id int) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
DELETE FROM users WHERE id=$1
`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

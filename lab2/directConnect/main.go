package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func connectDB() (*sql.DB, error) {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME_EASY")
	port := os.Getenv("DB_PORT")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected!")
	return db, nil
}

func createTable(db *sql.DB) {
	query := `
CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(50),
        age INT
    );
`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Unable to create table", err)
	}
	fmt.Println("Table created successfully")
}

func insertUser(db *sql.DB, name string, age int) {
	query := `
INSERT INTO users (name, age)
    VALUES ($1, $2)
`
	_, err := db.Exec(query, name, age)
	if err != nil {
		log.Fatalf("Unable to insert user", err)
	}
	fmt.Printf("Inserted user: %s, Age: %d\n", name, age)
}

func queryUsers(db *sql.DB) {
	query := `
SELECT id, name, age FROM users
`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Unable to query users", err)
	}
	fmt.Println("Users:")
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int
		err = rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatalf("Unable to scan row", err)
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalf("Error in rows", err)
	}
}

func main() {
	db, error := connectDB()
	if error != nil {
		log.Fatalf("Unable to connect to database", error)
	}
	defer db.Close()

	createTable(db)

	insertUser(db, "Asan", 19)
	insertUser(db, "Alisher", 20)
	insertUser(db, "Aibek", 21)

	queryUsers(db)

}

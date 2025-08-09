package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/go_events?parseTime=true"
	var err error

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening DB connection:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	fmt.Println("Database connected!")

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);
	`

	_, err := DB.Exec(createUsersTable)

	if err != nil {
		panic("Could not create users table")
	}

	createEventTables := `
		CREATE TABLE IF NOT EXISTS events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    location VARCHAR(255) NOT NULL,
    dateTime DATETIME NOT NULL,
    user_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

	`

	_, err = DB.Exec(createEventTables)

	if err != nil {
		panic("Could not create event tables")
	}

	createRegistrationTables := `
		CREATE TABLE IF NOT EXISTS events_registry (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT,
    user_id INT,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

	`

	_, err = DB.Exec(createRegistrationTables)

	if err != nil {
		panic("Could not create registrations tables")
	}
}

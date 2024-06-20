package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbSocket := os.Getenv("DB_SOCKET")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbSocket, dbName)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	createSchema(db)

	fmt.Println("Database initialized successfully")
}

func createSchema(db *sqlx.DB) {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL UNIQUE,
		email VARCHAR(255) NOT NULL UNIQUE,
		password CHAR(32) NOT NULL
	);`

	createBooksTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INT(11) AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		price DECIMAL(10,2) NOT NULL
	);`

	createOrdersTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT(11) NOT NULL,
		created_at datetime DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	createOrderBooksTable := `
	CREATE TABLE IF NOT EXISTS order_books (
		id INT AUTO_INCREMENT PRIMARY KEY,
		order_id INT(11),
		book_id INT(11),
		quantity INT(11),
		FOREIGN KEY (order_id) REFERENCES orders(id),
		FOREIGN KEY (book_id) REFERENCES books(id)
	);`

	db.MustExec(createUsersTable)
	db.MustExec(createBooksTable)
	db.MustExec(createOrdersTable)
	db.MustExec(createOrderBooksTable)
}

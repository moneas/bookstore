package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/moneas/bookstore/internal/domain/book"
	"github.com/moneas/bookstore/internal/domain/user"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Database configuration from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to the database
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Seed users
	seedUsers(db)

	// Seed books
	seedBooks(db)

	fmt.Println("Database seeded successfully")
}

func seedUsers(db *sqlx.DB) {
	users := []user.User{
		{Name: "John Doe", Email: "john@example.com", Username: "user1", Password: "password"},
		{Name: "Jane Smith", Email: "jane@example.com", Username: "user2", Password: "password"},
	}

	for _, u := range users {
		_, err := db.Exec("INSERT INTO users (name, email, username, password) VALUES (?, ?, ?, md5(?))", u.Name, u.Email, u.Username, u.Password)
		if err != nil {
			log.Fatalf("Error seeding users: %v", err)
		}
	}
}

func seedBooks(db *sqlx.DB) {
	books := []book.Book{
		{Title: "Go Programming", Author: "John Doe", Price: 29.99},
		{Title: "Learning Gin", Author: "Jane Smith", Price: 19.99},
	}

	for _, b := range books {
		_, err := db.Exec("INSERT INTO books (title, author, price) VALUES (?, ?, ?)", b.Title, b.Author, b.Price)
		if err != nil {
			log.Fatalf("Error seeding books: %v", err)
		}
	}
}

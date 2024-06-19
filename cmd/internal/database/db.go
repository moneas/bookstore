package database

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/moneas/bookstore/cmd/internal/models"
)

var DB *sqlx.DB

func Init() {
	dsn := "root:root@tcp(127.0.0.1:3306)/bookstore?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	schema := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE
    );

    CREATE TABLE IF NOT EXISTS books (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        author VARCHAR(255) NOT NULL,
        price DECIMAL(10,2) NOT NULL
    );

    CREATE TABLE IF NOT EXISTS orders (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );

    CREATE TABLE IF NOT EXISTS order_books (
        order_id INT,
        book_id INT,
        FOREIGN KEY (order_id) REFERENCES orders(id),
        FOREIGN KEY (book_id) REFERENCES books(id),
        PRIMARY KEY (order_id, book_id)
    );`

	DB.MustExec(schema)
}

func Seed() {
	books := []models.Book{
		{Title: "Book 1", Author: "Author 1", Price: 10.0},
		{Title: "Book 2", Author: "Author 2", Price: 15.0},
	}

	for _, book := range books {
		query := `INSERT INTO books (title, author, price) VALUES (?, ?, ?)`
		stmt, _ := DB.Prepare(query)
		stmt.Exec(book.Title, book.Author, book.Price)
		stmt.Close()
	}
}

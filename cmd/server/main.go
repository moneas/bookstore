package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/moneas/bookstore/internal/application"
	"github.com/moneas/bookstore/internal/infrastructure/database"
	"github.com/moneas/bookstore/internal/interfaces/http"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Get database configuration from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Initialize the database connection
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	defer db.Close()

	userRepository := database.NewUserRepository(db)
	userService := application.NewUserService(userRepository)

	orderRepository := database.NewOrderRepository(db)
	orderService := application.NewOrderService(orderRepository)

	bookRepository := database.NewBookRepository(db)
	bookService := application.NewBookService(bookRepository)

	// BIND user, order, and book handlers
	userHandler := http.NewUserHandler(userService)
	orderHandler := http.NewOrderHandler(orderService)
	bookHandler := http.NewBookHandler(bookService)

	r := gin.Default()

	// GIN routes
	r.GET("/books", bookHandler.GetBooks)
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.GetUsers)
	r.POST("/orders", http.Authenticate(userService), orderHandler.CreateOrder)
	r.GET("/myorders", http.Authenticate(userService), orderHandler.GetOrdersByUserID)

	// run GIN on port 8080
	r.Run(":8080")
}

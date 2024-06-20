// cmd/server/main.go

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
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbSocket := os.Getenv("DB_SOCKET")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbSocket, dbName)

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

	// BIND user and order handlers
	userHandler := http.NewUserHandler(userService)
	orderHandler := http.NewOrderHandler(orderService)

	r := gin.Default()

	// GIN routes
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.GetUsers)
	r.POST("/orders", orderHandler.CreateOrder)
	r.GET("/ordersbyuser/:user_id", orderHandler.GetOrdersByUserID)

	// run GIN on port 8080
	r.Run(":8080")
}

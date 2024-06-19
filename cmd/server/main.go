package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moneas/bookstore/internal/database"
	"github.com/moneas/bookstore/internal/handlers"
)

func main() {
	database.Init()
	database.Seed()

	r := gin.Default()

	r.POST("/users", handlers.CreateUser)
	r.GET("/books", handlers.GetBooks)
	r.POST("/orders", handlers.CreateOrder)
	r.GET("/orders", handlers.GetOrders)

	r.Run(":8080")
}

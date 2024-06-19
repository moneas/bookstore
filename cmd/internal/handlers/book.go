package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moneas/bookstore/internal/database"
	"github.com/moneas/bookstore/internal/models"
)

func GetBooks(c *gin.Context) {
	var books []models.Book
	query := `SELECT * FROM books`
	err := database.DB.Select(&books, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

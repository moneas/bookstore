package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moneas/bookstore/internal/database"
	"github.com/moneas/bookstore/internal/models"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO users (name, email) VALUES (?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Name, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	user.ID = uint(id)
	c.JSON(http.StatusCreated, user)
}

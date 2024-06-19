package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moneas/bookstore/internal/database"
	"github.com/moneas/bookstore/internal/models"
)

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.MustBegin()

	query := `INSERT INTO orders (user_id) VALUES (?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(order.UserID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderID, _ := res.LastInsertId()
	order.ID = uint(orderID)

	for _, book := range order.Books {
		stmt, err := tx.Prepare(`INSERT INTO order_books (order_id, book_id) VALUES (?, ?)`)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(orderID, book.ID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, order)
}

func GetOrders(c *gin.Context) {
	var orders []models.Order
	query := `SELECT * FROM orders`
	err := database.DB.Select(&orders, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i, order := range orders {
		var books []models.Book
		query := `SELECT b.* FROM books b INNER JOIN order_books ob ON b.id = ob.book_id WHERE ob.order_id = ?`
		err := database.DB.Select(&books, query, order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orders[i].Books = books
	}

	c.JSON(http.StatusOK, orders)
}

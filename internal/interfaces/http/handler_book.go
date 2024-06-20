package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moneas/bookstore/internal/application"
	"github.com/moneas/bookstore/internal/domain/book"
)

type BookHandler struct {
	bookService *application.BookService
}

func NewBookHandler(bookService *application.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	books, err := h.bookService.GetBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var book book.Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book data"})
		return
	}

	createdBook, err := h.bookService.CreateBook(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}
	c.JSON(http.StatusCreated, createdBook)
}

func (h *BookHandler) GetBookByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book, err := h.bookService.GetBookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

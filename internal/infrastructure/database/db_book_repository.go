package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/moneas/bookstore/internal/domain/book"
)

type BookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) FindAll() ([]book.Book, error) {
	var books []book.Book
	err := r.db.Select(&books, "SELECT * FROM books")
	return books, err
}

func (r *BookRepository) FindByID(id uint) (*book.Book, error) {
	var book book.Book
	err := r.db.Get(&book, "SELECT * FROM books WHERE id=?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) Save(book *book.Book) error {
	_, err := r.db.NamedExec("INSERT INTO books (title, author, price) VALUES (:title, :author, :price)", book)
	return err
}

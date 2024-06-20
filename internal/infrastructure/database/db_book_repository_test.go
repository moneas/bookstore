package database

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/moneas/bookstore/internal/domain/book"
	"github.com/stretchr/testify/assert"
)

func TestFindAllBooks_Positive(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewBookRepository(sqlxDB)

	mock.ExpectQuery(`SELECT \* FROM books`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "price"}).
			AddRow(1, "Book 1", "Author 1", 10.0).
			AddRow(2, "Book 2", "Author 2", 15.0))

	books, err := repo.FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, books)
	assert.Equal(t, 2, len(books))
	assert.Equal(t, uint(1), books[0].ID)
	assert.Equal(t, "Book 1", books[0].Title)
	assert.Equal(t, "Author 1", books[0].Author)
	assert.Equal(t, 10.0, books[0].Price)
}

func TestFindAllBooks_Negative_ErrorQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewBookRepository(sqlxDB)

	mock.ExpectQuery(`SELECT \* FROM books`).
		WillReturnError(fmt.Errorf("query error"))

	books, err := repo.FindAll()
	assert.Error(t, err)
	assert.Nil(t, books)
}

func TestFindBookByID_Positive(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewBookRepository(sqlxDB)

	mock.ExpectQuery(`SELECT \* FROM books WHERE id=\?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "price"}).
			AddRow(1, "Book 1", "Author 1", 10.0))

	book, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, uint(1), book.ID)
	assert.Equal(t, "Book 1", book.Title)
	assert.Equal(t, "Author 1", book.Author)
	assert.Equal(t, 10.0, book.Price)
}

func TestFindBookByID_Negative_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewBookRepository(sqlxDB)

	mock.ExpectQuery(`SELECT \* FROM books WHERE id=\?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "price"}))

	book, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.Nil(t, book)
}

func TestFindBookByID_Negative_ErrorQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewBookRepository(sqlxDB)

	mock.ExpectQuery(`SELECT \* FROM books WHERE id=\?`).
		WithArgs(1).
		WillReturnError(fmt.Errorf("query error"))

	book, err := repo.FindByID(1)
	assert.Error(t, err)
	assert.Nil(t, book)
}

func TestSaveBook_Positive(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewBookRepository(sqlxDB)

	mock.ExpectExec(`INSERT INTO books \(title, author, price\) VALUES \(\?, \?, \?\)`).
		WithArgs("Book 1", "Author 1", 10.0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	newBook := &book.Book{
		Title:  "Book 1",
		Author: "Author 1",
		Price:  10.0,
	}
	err = repo.Save(newBook)
	assert.NoError(t, err)
}

func TestSaveBook_Negative_ErrorQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewBookRepository(sqlxDB)

	mock.ExpectExec(`INSERT INTO books \(title, author, price\) VALUES \(\?, \?, \?\)`).
		WithArgs("Book 1", "Author 1", 10.0).
		WillReturnError(fmt.Errorf("insert error"))

	newBook := &book.Book{
		Title:  "Book 1",
		Author: "Author 1",
		Price:  10.0,
	}
	err = repo.Save(newBook)
	assert.Error(t, err)
}

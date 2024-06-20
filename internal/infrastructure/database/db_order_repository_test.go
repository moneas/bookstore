package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func TestFindOrderDetails_Positive(t *testing.T) {
	db, mock := newTestDB(t)
	repo := NewOrderRepository(db)

	orderID := 1
	mock.ExpectQuery(`
		SELECT
			ob.quantity as quantity,
			b.id AS book_id,
			b.title as title,
			b.author,
			b.price
		FROM orders o
		INNER JOIN order_books ob ON o.id = ob.order_id
		INNER JOIN books b ON ob.book_id = b.id
		WHERE o.id = ?`).
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"quantity", "book_id", "title", "author", "price"}).
			AddRow(2, 1, "Book Title", "Book Author", 19.99))

	orderDetail, err := repo.FindOrderDetails(orderID)
	assert.NoError(t, err)
	assert.NotNil(t, orderDetail)
	assert.Equal(t, 1, len(orderDetail.Books))
	assert.Equal(t, "Book Title", orderDetail.Books[0].Title)
}

func TestFindOrderDetails_Negative_NoBooks(t *testing.T) {
	db, mock := newTestDB(t)
	repo := NewOrderRepository(db)

	orderID := 1
	mock.ExpectQuery(`
		SELECT
			ob.quantity as quantity,
			b.id AS book_id,
			b.title as title,
			b.author,
			b.price
		FROM orders o
		INNER JOIN order_books ob ON o.id = ob.order_id
		INNER JOIN books b ON ob.book_id = b.id
		WHERE o.id = ?`).
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"quantity", "book_id", "title", "author", "price"}))

	orderDetail, err := repo.FindOrderDetails(orderID)
	assert.NoError(t, err)
	assert.NotNil(t, orderDetail)
	assert.Equal(t, 0, len(orderDetail.Books))
}

func TestFindOrderDetails_Negative_ErrorQuery(t *testing.T) {
	db, mock := newTestDB(t)
	repo := NewOrderRepository(db)

	orderID := 1
	mock.ExpectQuery(`
		SELECT
			ob.quantity as quantity,
			b.id AS book_id,
			b.title as title,
			b.author,
			b.price
		FROM orders o
		INNER JOIN order_books ob ON o.id = ob.order_id
		INNER JOIN books b ON ob.book_id = b.id
		WHERE o.id = ?`).
		WithArgs(orderID).
		WillReturnError(fmt.Errorf("query error"))

	orderDetail, err := repo.FindOrderDetails(orderID)
	assert.Error(t, err)
	assert.Nil(t, orderDetail)
}

func TestFindOrdersByUserID_Positive(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewOrderRepository(sqlxDB)

	createdAt := time.Now()

	mock.ExpectQuery(`SELECT id, user_id, created_at FROM orders WHERE user_id = \?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "created_at"}).
			AddRow(1, 1, createdAt))

	mock.ExpectQuery(`SELECT b.id as book_id, b.title, b.author, b.price, ob.quantity FROM books b INNER JOIN order_books ob ON b.id = ob.book_id WHERE ob.order_id = \?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"book_id", "title", "author", "price", "quantity"}).
			AddRow(1, "Book 1", "Author 1", 10.0, 2).
			AddRow(2, "Book 2", "Author 2", 15.0, 1))

	orders, err := repo.FindOrdersByUserID(1)
	assert.NoError(t, err)
	assert.NotNil(t, orders)
	assert.Equal(t, 1, len(orders))
	assert.Equal(t, uint(1), orders[0].ID)
	assert.Equal(t, uint(1), orders[0].UserID)
	assert.Equal(t, createdAt, orders[0].CreatedAt)
	assert.Equal(t, 2, len(orders[0].Books))
	assert.Equal(t, uint(1), orders[0].Books[0].ID)
	assert.Equal(t, "Book 1", orders[0].Books[0].Title)
	assert.Equal(t, 2, orders[0].Books[0].Quantity)
}

func TestFindOrdersByUserID_Negative_ErrorQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewOrderRepository(sqlxDB)

	mock.ExpectQuery(`SELECT id, user_id, created_at FROM orders WHERE user_id = \?`).
		WithArgs(1).
		WillReturnError(fmt.Errorf("query error"))

	orders, err := repo.FindOrdersByUserID(1)
	assert.Error(t, err)
	assert.Nil(t, orders)
}

func TestFindOrdersByUserID_Negative_InvalidScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewOrderRepository(sqlxDB)

	mock.ExpectQuery(`SELECT id, user_id, created_at FROM orders WHERE user_id = \?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "created_at"}).
			AddRow(1, "invalid_user_id", "invalid_created_at"))

	orders, err := repo.FindOrdersByUserID(1)
	assert.Error(t, err)
	assert.Nil(t, orders)
}

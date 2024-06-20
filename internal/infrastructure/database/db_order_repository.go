package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/moneas/bookstore/internal/domain/order"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) FindAll() ([]order.Order, error) {
	var orders []order.Order
	err := r.db.Select(&orders, "SELECT * FROM orders")
	return orders, err
}

func (r *OrderRepository) FindByID(id uint) (*order.Order, error) {
	var order order.Order
	err := r.db.Get(&order, "SELECT * FROM orders WHERE id=?", id)
	return &order, err
}

func (r *OrderRepository) Save(order *order.Order) error {
	// Start a transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert into orders table
	result, err := tx.NamedExec("INSERT INTO orders (user_id) VALUES (:user_id)", order)
	if err != nil {
		return err
	}

	// Get the last inserted ID
	orderID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	order.ID = uint(orderID)

	// Insert into order_books table
	for _, book := range order.Books {
		_, err := tx.Exec("INSERT INTO order_books (order_id, book_id, quantity) VALUES (?, ?, ?)", order.ID, book.ID, book.Quantity)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) FindOrderDetails(orderID int) (*order.OrderDetail, error) {
	query := `
		SELECT
			ob.quantity as quantity,
			b.id AS book_id,
			b.title as title,
			b.author,
			b.price
		FROM orders o
		INNER JOIN order_books ob ON o.id = ob.order_id
		INNER JOIN books b ON ob.book_id = b.id
		WHERE o.id = ?
	`
	var userOrder order.OrderDetail
	var books []order.BookWithQuantity

	defaultOrder, _ := r.FindByID(uint(orderID))

	rows, err := r.db.Queryx(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map rows to order struct
	for rows.Next() {
		var book order.BookWithQuantity
		err := rows.StructScan(&book)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	userOrder.ID = uint(orderID)
	userOrder.Books = books
	userOrder.UserID = defaultOrder.UserID

	return &userOrder, nil
}

func (r *OrderRepository) FindOrdersByUserID(userID int) ([]*order.OrderDetail, error) {
	var orders []*order.OrderDetail

	query := `SELECT id, user_id, created_at FROM orders WHERE user_id = ?`
	err := r.db.Select(&orders, query, userID)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		//no order found
		return orders, nil
	}

	for _, orderDetail := range orders {
		query = `
            SELECT b.id as book_id, b.title, b.author, b.price, ob.quantity
            FROM books b
            INNER JOIN order_books ob ON b.id = ob.book_id
            WHERE ob.order_id = ?`
		rows, err := r.db.Queryx(query, orderDetail.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var books []order.BookWithQuantity
		for rows.Next() {
			var book order.BookWithQuantity
			err := rows.StructScan(&book)
			if err != nil {
				return nil, err
			}
			books = append(books, book)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		orderDetail.Books = books
	}

	return orders, nil
}

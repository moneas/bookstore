package order

import (
	"time"
)

type Order struct {
	ID        uint        `db:"id" json:"id"`
	UserID    uint        `db:"user_id" json:"user_id"`
	Books     []OrderBook `json:"books"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
}

type OrderBook struct {
	ID       uint `db:"id" json:"id"`
	Quantity uint `db:"quantity" json:"quantity"`
}

type BookWithQuantity struct {
	ID       uint    `db:"book_id" json:"id"`
	Title    string  `db:"title" json:"title"`
	Author   string  `db:"author" json:"author"`
	Price    float64 `db:"price" json:"price"`
	Quantity int     `db:"quantity" json:"quantity"`
}

type OrderDetail struct {
	ID        uint               `db:"id" json:"id"`
	UserID    uint               `db:"user_id" json:"user_id"`
	Books     []BookWithQuantity `json:"books"`
	CreatedAt time.Time          `db:"created_at" json:"created_at"`
}

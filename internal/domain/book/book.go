package book

type Book struct {
	ID     uint    `db:"id" json:"id"`
	Title  string  `db:"title" json:"title"`
	Author string  `db:"author" json:"author"`
	Price  float64 `db:"price" json:"price"`
}

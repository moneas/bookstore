package models

type Order struct {
	ID     uint   `db:"id" json:"id"`
	UserID uint   `db:"user_id" json:"user_id"`
	Books  []Book `json:"books"`
}

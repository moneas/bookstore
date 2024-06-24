package user

type User struct {
	ID       uint   `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

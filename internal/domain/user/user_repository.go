package user

type Repository interface {
	FindAll() ([]User, error)
	FindByID(id uint) (*User, error)
	Save(user *User) error
	FindByEmail(email string) (*User, error)
}

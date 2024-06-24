package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/moneas/bookstore/internal/domain/user"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll() ([]user.User, error) {
	var users []user.User
	err := r.db.Select(&users, "SELECT * FROM users")
	return users, err
}

func (r *UserRepository) FindByID(id uint) (*user.User, error) {
	var user user.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id=?", id)
	return &user, err
}

func (r *UserRepository) Save(user *user.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.NamedExec(`
        INSERT INTO users (name, email, username, password)
        VALUES (:name, :email, :username, :password)
    `, user)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint(id)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	var user user.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE email=?", email)
	return &user, err
}

func (r *UserRepository) FindByUsername(username string) (*user.User, error) {
	var user user.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE username=?", username)
	return &user, err
}

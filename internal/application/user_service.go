// application/user_service.go
package application

import (
	"errors"

	"github.com/moneas/bookstore/internal/domain/user"
	"github.com/moneas/bookstore/internal/util"
)

type UserService struct {
	userRepository user.Repository
}

func NewUserService(userRepository user.Repository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUsers() ([]user.User, error) {
	return s.userRepository.FindAll()
}

func (s *UserService) GetUserByID(id uint) (*user.User, error) {
	return s.userRepository.FindByID(id)
}

func (s *UserService) CreateUser(user *user.User) error {
	return s.userRepository.Save(user)
}

func (s *UserService) GetUserByEmail(email string) (*user.User, error) {
	return s.userRepository.FindByEmail(email)
}

func (s *UserService) Authenticate(username, password string) (*user.User, error) {
	u, err := s.userRepository.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !util.CheckPasswordHash(password, u.Password) {
		return nil, errors.New("incorrect password")
	}

	return u, nil
}

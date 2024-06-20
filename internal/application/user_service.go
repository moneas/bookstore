package application

import (
	"github.com/moneas/bookstore/internal/domain/user"
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

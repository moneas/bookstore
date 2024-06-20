package application

import (
	"github.com/moneas/bookstore/internal/domain/book"
)

type BookService struct {
	repo book.Repository
}

func NewBookService(repo book.Repository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) GetBooks() ([]book.Book, error) {
	return s.repo.FindAll()
}

func (s *BookService) GetBookByID(id uint) (*book.Book, error) {
	return s.repo.FindByID(id)
}

func (s *BookService) CreateBook(book *book.Book) (*book.Book, error) {
	if err := s.repo.Save(book); err != nil {
		return nil, err
	}
	return book, nil
}

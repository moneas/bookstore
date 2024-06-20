package book

type Repository interface {
	FindAll() ([]Book, error)
	FindByID(id uint) (*Book, error)
	Save(book *Book) error
}

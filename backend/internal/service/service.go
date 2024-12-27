package service

import (
	"cmd/main.go/internal/storage"
	"cmd/main.go/models"
)

type Service interface {
	GetBooks() ([]models.Book, error)
	AddBook(book models.Book) (models.Book, error)
	DeleteBook(id int) error

	GetUsers() ([]models.User, error)
	AddUser(user models.User) (models.User, error)
	DeleteUser(id int) error

	GetLoans() ([]models.Loan, error)
	IssueLoan(userID, bookID int) (models.Loan, error)
	ReturnLoan(loanID int) (models.Loan, error)

	GetStatistics() (models.Statistics, error)
}

type service struct {
	db storage.Database
}

func (s service) GetBooks() ([]models.Book, error) {
	return s.db.GetBooks()
}

func (s service) AddBook(book models.Book) (models.Book, error) {
	id, err := s.db.AddBook(book)
	if err != nil {
		return models.Book{}, err
	}
	book.ID = id
	return book, nil
}

func (s service) DeleteBook(id int) error {
	return s.db.DeleteBook(id)
}

func (s service) GetUsers() ([]models.User, error) {
	return s.db.GetUsers()
}

func (s service) AddUser(user models.User) (models.User, error) {
	return s.db.AddUser(user)
}

func (s service) DeleteUser(id int) error {
	return s.db.DeleteUser(id)
}

func (s service) GetLoans() ([]models.Loan, error) {
	return s.db.GetLoans()
}

func (s service) IssueLoan(userID, bookID int) (models.Loan, error) {
	return s.db.IssueLoan(userID, bookID)
}

func (s service) ReturnLoan(loanID int) (models.Loan, error) {
	return s.db.ReturnLoan(loanID)
}

func (s service) GetStatistics() (models.Statistics, error) {
	stats, err := s.db.GetStats()
	return *stats, err
}

// NewService создает новый экземпляр сервиса
func NewService(db storage.Database) Service {
	return service{
		db: db,
	}
}

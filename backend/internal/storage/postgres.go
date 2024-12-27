package storage

import (
	"cmd/main.go/models"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Подключаем драйвер PostgreSQL
	"strconv"
	"time"
)

type Database struct {
	db *sql.DB
}

// NewDatabase инициализирует подключение к базе данных PostgreSQL.
func NewDatabase(dataSourceName string) Database {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return Database{}
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("Error pinging database:", err))
		return Database{}
	}

	return Database{db: db}
}

// Close закрывает подключение к базе данных.
func (d *Database) Close() error {
	return d.db.Close()
}

// Migrate применяет схему для приложения.
func (d *Database) Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			category TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE
		);`,
		`CREATE TABLE IF NOT EXISTS loans (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			book_id INT NOT NULL,
			borrow_date DATE NOT NULL,
			return_date DATE,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (book_id) REFERENCES books(id)
		);`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

// CRUD операции для книг

func (d *Database) GetBooks() ([]models.Book, error) {
	rows, err := d.db.Query("SELECT id, title, author, category FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Category); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (d *Database) AddBook(book models.Book) (int, error) {
	var id int
	err := d.db.QueryRow("INSERT INTO books (title, author, category) VALUES ($1, $2, $3) RETURNING id", book.Title, book.Author, book.Category).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *Database) DeleteBook(id int) error {
	_, err := d.db.Exec("DELETE FROM books WHERE id = $1", id)
	return err
}

// CRUD операции для пользователей

func (d *Database) GetUsers() ([]models.User, error) {
	rows, err := d.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (d *Database) AddUser(user models.User) (models.User, error) {
	var id int
	err := d.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&id)
	if err != nil {
		return models.User{}, err
	}

	user.ID = id
	return user, nil
}

func (d *Database) DeleteUser(id int) error {
	_, err := d.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// CRUD операции для займов
func (d *Database) GetLoans() ([]models.Loan, error) {
	rows, err := d.db.Query("SELECT id, user_id, book_id, borrow_date, return_date FROM loans")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []models.Loan
	for rows.Next() {
		var loan models.Loan
		// Используем time.Time для полей даты
		if err := rows.Scan(&loan.ID, &loan.UserID, &loan.BookID, &loan.BorrowDate, &loan.ReturnDate); err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	return loans, nil
}

func (d *Database) IssueLoan(userID, bookID int) (models.Loan, error) {
	var id int
	borrowDate := time.Now() // Используем текущую дату и время
	err := d.db.QueryRow("INSERT INTO loans (user_id, book_id, borrow_date) VALUES ($1, $2, $3) RETURNING id", userID, bookID, borrowDate).Scan(&id)
	if err != nil {
		return models.Loan{}, err
	}

	// Возвращаем структуру с типом time.Time для даты
	return models.Loan{
		ID:         id,
		UserID:     strconv.Itoa(userID),
		BookID:     strconv.Itoa(bookID),
		BorrowDate: borrowDate.String(), // Используем тип time.Time
	}, nil
}

func (d *Database) ReturnLoan(loanID int) (models.Loan, error) {
	// Обновляем returnDate с использованием текущей даты
	returnDate := time.Now() // Используем time.Now(), чтобы получить текущую дату
	_, err := d.db.Exec("UPDATE loans SET return_date = $1 WHERE id = $2", returnDate, loanID)
	if err != nil {
		return models.Loan{}, err
	}
	var date string = returnDate.String()
	// Возвращаем структуру с правильным типом для ReturnDate
	return models.Loan{
		ID:         loanID,
		ReturnDate: &date, // Используем указатель на time.Time для nullable даты
	}, nil
}

// GetStats собирает статистику по библиотеке
func (r *Database) GetStats() (*models.Statistics, error) {
	stats := &models.Statistics{}

	// Получаем общее количество книг
	err := r.db.QueryRow("SELECT COUNT(*) FROM books").Scan(&stats.TotalBooks)
	if err != nil {
		return nil, fmt.Errorf("error fetching total books: %v", err)
	}

	// Получаем общее количество пользователей
	err = r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	if err != nil {
		return nil, fmt.Errorf("error fetching total users: %v", err)
	}

	// Получаем общее количество займов
	err = r.db.QueryRow("SELECT COUNT(*) FROM loans WHERE return_date IS NULL").Scan(&stats.TotalLoans)
	if err != nil {
		return nil, fmt.Errorf("error fetching total loans: %v", err)
	}

	// Получаем статистику по категориям
	rows, err := r.db.Query(`
		SELECT category, COUNT(loans.id) 
		FROM books
		LEFT JOIN loans ON books.id = loans.book_id 
		WHERE loans.return_date IS NULL
		GROUP BY category`)
	if err != nil {
		return nil, fmt.Errorf("error fetching category stats: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var categoryStats models.CategoryStats
		err := rows.Scan(&categoryStats.Name, &categoryStats.Count)
		if err != nil {
			return nil, fmt.Errorf("error scanning category stats: %v", err)
		}
		stats.PopularCategories = append(stats.PopularCategories, categoryStats)
	}

	// Получаем статистику по пользователям
	rows, err = r.db.Query(`
		SELECT users.name, COUNT(loans.id)
		FROM users
		LEFT JOIN loans ON users.id = loans.user_id
		WHERE loans.return_date IS NULL
		GROUP BY users.name`)
	if err != nil {
		return nil, fmt.Errorf("error fetching user stats: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userStats models.UserStats
		err := rows.Scan(&userStats.Name, &userStats.LoansCount)
		if err != nil {
			return nil, fmt.Errorf("error scanning user stats: %v", err)
		}
		stats.ActiveUsers = append(stats.ActiveUsers, userStats)
	}

	// Возвращаем собранную статистику
	return stats, nil
}

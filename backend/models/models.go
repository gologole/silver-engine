package models

// Book represents a book in the library system
type Book struct {
	ID       int
	Title    string
	Author   string
	Category string
}

// User represents a user in the library system
type User struct {
	ID            int
	Name          string
	Email         string
	BorrowedBooks []string
}

// Loan represents a book loan
type Loan struct {
	ID         int     `json:"id"`
	UserID     string  `json:"userID"`
	BookID     string  `json:"bookID"`
	BorrowDate string  `json:"borrowDate"`
	ReturnDate *string `json:"returnDate"`
}

// Statistics represents library statistics
type Statistics struct {
	TotalBooks        int
	TotalUsers        int
	TotalLoans        int
	PopularCategories []CategoryStats
	ActiveUsers       []UserStats
}

// CategoryStats represents statistics for a category
type CategoryStats struct {
	Name  string
	Count int //количество loans в категории
}

// UserStats represents statistics for a user
type UserStats struct {
	Name       string
	LoansCount int //количество loans пользователя
}

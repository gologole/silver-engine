package transport

import (
	"cmd/main.go/internal/service"
	"cmd/main.go/models"
	"github.com/gin-contrib/cors" // Импортируем пакет
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Handler struct для работы с сервисом
type Handler struct {
	service service.Service
}

// NewHandler создает новый обработчик с сервисом
func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

// InitRoutes инициализирует маршруты для обработки HTTP запросов
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	// Добавляем CORS middleware
	router.Use(cors.Default()) // Вы можете настроить CORS здесь, если нужно

	api := router.Group("/api")

	// Books
	books := api.Group("/books")
	{
		books.GET("", h.GetBooks)
		books.POST("", h.AddBook)
		books.DELETE(":id", h.DeleteBook)
	}

	// Users
	users := api.Group("/users")
	{
		users.GET("", h.GetUsers)
		users.POST("", h.AddUser)
		users.DELETE(":id", h.DeleteUser)
	}

	// Loans
	loans := api.Group("/loans")
	{
		loans.GET("", h.GetLoans)
		loans.POST("", h.IssueLoan)
		loans.POST(":id/return", h.ReturnLoan)
	}

	// Statistics
	api.GET("/statistics", h.GetStatistics)
	return router
}

// Books Handlers

// GetBooks обрабатывает запрос на получение списка книг
func (h *Handler) GetBooks(c *gin.Context) {
	books, err := h.service.GetBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

// AddBook обрабатывает запрос на добавление книги
func (h *Handler) AddBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBook, err := h.service.AddBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newBook)
}

// DeleteBook обрабатывает запрос на удаление книги
func (h *Handler) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := h.service.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

// Users Handlers

// GetUsers обрабатывает запрос на получение списка пользователей
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// AddUser обрабатывает запрос на добавление нового пользователя
func (h *Handler) AddUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := h.service.AddUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// DeleteUser обрабатывает запрос на удаление пользователя
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// Loans Handlers

// GetLoans обрабатывает запрос на получение списка всех займов
func (h *Handler) GetLoans(c *gin.Context) {
	loans, err := h.service.GetLoans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loans)
}

// IssueLoan обрабатывает запрос на выдачу книги
func (h *Handler) IssueLoan(c *gin.Context) {
	var loan models.Loan
	if err := c.ShouldBindJSON(&loan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Преобразуем строковый userID и bookID в целые числа, если нужно
	userID, err := strconv.Atoi(loan.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
		return
	}

	bookID, err := strconv.Atoi(loan.BookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookID"})
		return
	}

	newLoan, err := h.service.IssueLoan(userID, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusCreated, newLoan)
}

// ReturnLoan обрабатывает запрос на возврат книги
func (h *Handler) ReturnLoan(c *gin.Context) {
	loanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	updatedLoan, err := h.service.ReturnLoan(loanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedLoan)
}

// Statistics Handler

// GetStatistics обрабатывает запрос на получение статистики по библиотеке
func (h *Handler) GetStatistics(c *gin.Context) {
	stats, err := h.service.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

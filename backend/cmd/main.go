package main

import (
	config2 "cmd/main.go/config"
	"cmd/main.go/internal/service"
	"cmd/main.go/internal/storage"
	"cmd/main.go/internal/transport"
	"cmd/main.go/pkg/logger"
	"cmd/main.go/server"

	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Load configuration
	config, err := config2.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Error loading configuration: %s", err.Error()))
	}
	// Initialize logger
	appLogger := logger.NewLogger()
	defer appLogger.Sync()

	// Create context with logger
	ctx := logger.WithLoggerContext(context.Background(), appLogger)

	// Connect to PostgreSQL
	db := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	appLogger.Info("Connecting to PostgreSQL", zap.String("db_url", db))

	mystorage := storage.NewDatabase(db)
	defer mystorage.Close()
	mystorage.Migrate()
	myservice := service.NewService(mystorage)
	myhandler := transport.NewHandler(myservice)

	// Initialize server
	srv := &server.Server{}

	go func() {
		if err := srv.RunServer(config.AppPort, myhandler.InitRoutes()); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	appLogger.Info("Server is running", zap.String("port", config.AppPort))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server shutdown failed", zap.Error(err))
	}

	mystorage.Close()
	appLogger.Info("Server shutdown successfully")
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"rutube/controller"
	telegramconnect "rutube/infrastructure/TelegramConnect"
	"rutube/infrastructure/database"
	"rutube/infrastructure/router"
	"rutube/infrastructure/server"
	"rutube/usecase"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Error starting logger. detailed information - ", err.Error())
		os.Exit(1)
	}

	db, err := database.InitDatabase("Date.db")
	if err != nil {
		logger.Error("Database initialization error", zap.Error(err))
		os.Exit(1)
	}
	defer db.Close()

	tg, err := telegramconnect.NewTelegramClient(logger)
	if err != nil {
		logger.Error("Telegram client initialization error", zap.Error(err))
		os.Exit(1)
	}

	dbService := database.NewDatabase(logger, db)
	useCase := usecase.NewUseCase(logger, dbService, tg)
	handler := controller.NewHandlers(logger, useCase)
	rtr := router.NewGoChiRouting(logger, handler)
	srv := server.NewServerHTTP(logger, rtr, ":8080")

	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("Error starting the server", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	logger.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}

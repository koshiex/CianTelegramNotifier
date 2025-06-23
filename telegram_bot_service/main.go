package main

import (
	"log"
	"os"
	"telegram_bot_service/internal/bot"
	"telegram_bot_service/internal/config"
	"telegram_bot_service/internal/database"
	"telegram_bot_service/internal/services"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found")
	}

	// Initialize config
	cfg := config.New()

	// Setup logging
	setupLogging(cfg.LogLevel)

	// Initialize database
	db, err := database.Initialize(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize services
	cianService := services.NewCianService(cfg.CianAPIURL)
	userService := services.NewUserService(db)
	favoriteService := services.NewFavoriteService(db)

	// Initialize and start bot
	telegramBot, err := bot.New(cfg.TelegramToken, cianService, userService, favoriteService)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Start health check server if enabled
	if cfg.HealthCheckEnabled {
		healthServer := bot.NewHealthServer(telegramBot, cfg.HealthCheckPort)
		healthServer.Start()
	}

	logrus.Info("Starting Telegram bot...")
	if err := telegramBot.Start(); err != nil {
		log.Fatalf("Bot error: %v", err)
	}
}

func setupLogging(level string) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetOutput(os.Stdout)
}

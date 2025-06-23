package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type HealthServer struct {
	bot  *Bot
	port string
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Services  struct {
		TelegramBot bool `json:"telegram_bot"`
		CianAPI     bool `json:"cian_api"`
		Database    bool `json:"database"`
	} `json:"services"`
}

func NewHealthServer(bot *Bot, port string) *HealthServer {
	return &HealthServer{
		bot:  bot,
		port: port,
	}
}

func (hs *HealthServer) Start() {
	http.HandleFunc("/health", hs.healthHandler)
	http.HandleFunc("/ready", hs.readinessHandler)

	logrus.WithField("port", hs.port).Info("Starting health check server")

	go func() {
		if err := http.ListenAndServe(":"+hs.port, nil); err != nil {
			logrus.WithError(err).Error("Health server failed")
		}
	}()
}

func (hs *HealthServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	// Check Telegram Bot
	response.Services.TelegramBot = hs.bot.api != nil

	// Check CIAN API
	if err := hs.bot.cianService.HealthCheck(); err != nil {
		response.Services.CianAPI = false
		response.Status = "degraded"
		logrus.WithError(err).Warn("CIAN API health check failed")
	} else {
		response.Services.CianAPI = true
	}

	// Check Database (simple check)
	if _, err := hs.bot.userService.GetAllActiveUsers(); err != nil {
		response.Services.Database = false
		response.Status = "degraded"
		logrus.WithError(err).Warn("Database health check failed")
	} else {
		response.Services.Database = true
	}

	w.Header().Set("Content-Type", "application/json")

	statusCode := http.StatusOK
	if response.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (hs *HealthServer) readinessHandler(w http.ResponseWriter, r *http.Request) {
	// Simple readiness check
	if hs.bot.api == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Bot not ready")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ready")
}

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"telegram_bot_service/internal/models"
	"time"

	"github.com/sirupsen/logrus"
)

type CianService struct {
	baseURL    string
	httpClient *http.Client
}

func NewCianService(baseURL string) *CianService {
	return &CianService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// GetListings fetches listings from CIAN API
func (s *CianService) GetListings(forceRefresh bool) ([]models.Listing, error) {
	url := fmt.Sprintf("%s/listings", s.baseURL)
	if forceRefresh {
		url += "?refresh=true"
	}

	resp, err := s.httpClient.Get(url)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch listings")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var listings []models.Listing
	if err := json.Unmarshal(body, &listings); err != nil {
		return nil, err
	}

	logrus.WithField("count", len(listings)).Debug("Fetched listings from CIAN API")
	return listings, nil
}

// GetSettings fetches current search settings
func (s *CianService) GetSettings() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/settings", s.baseURL)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch settings")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, err
	}

	logrus.Debug("Fetched settings from CIAN API")
	return settings, nil
}

// UpdateSettings updates search settings
func (s *CianService) UpdateSettings(settings map[string]interface{}) error {
	url := fmt.Sprintf("%s/settings", s.baseURL)

	jsonData, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to update settings")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	logrus.Debug("Updated settings via CIAN API")
	return nil
}

// HealthCheck checks if CIAN API is healthy
func (s *CianService) HealthCheck() error {
	url := fmt.Sprintf("%s/health", s.baseURL)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}

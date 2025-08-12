package services

import (
	"fmt"

	"github.com/adeesh/url-shortener/internal/constants"
	"github.com/adeesh/url-shortener/internal/database"
)

// AnalyticsService handles analytics tracking.
type AnalyticsService struct{}

// NewAnalyticsService creates a new analytics service instance.
func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

// TrackRedirectCounter increments the total redirect counter.
func (s *AnalyticsService) TrackRedirectCounter() error {
	r := database.CreateClient(constants.RedisDBRateLimit)
	defer func() {
		if err := database.CloseClient(r); err != nil {
			// Log error but don't fail the main operation
			_ = err
		}
	}()

	return database.Increment(r, constants.Counter)
}

// GetRedirectCount returns the total number of redirects.
func (s *AnalyticsService) GetRedirectCount() (int64, error) {
	r := database.CreateClient(constants.RedisDBRateLimit)
	defer func() {
		if err := database.CloseClient(r); err != nil {
			// Log error but don't fail the main operation
			_ = err
		}
	}()

	val, err := database.Get(r, constants.Counter)
	if err != nil {
		return 0, err
	}

	// Parse the string value to int64
	var count int64
	_, err = fmt.Sscanf(val, "%d", &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// TrackShortURLAccess tracks access to a specific short URL.
func (s *AnalyticsService) TrackShortURLAccess(shortCode string) error {
	r := database.CreateClient(constants.RedisDBRateLimit)
	defer func() {
		if err := database.CloseClient(r); err != nil {
			// Log error but don't fail the main operation
			_ = err
		}
	}()

	key := "access:" + shortCode
	return database.Increment(r, key)
}

// GetShortURLAccessCount returns the access count for a specific short URL.
func (s *AnalyticsService) GetShortURLAccessCount(shortCode string) (int64, error) {
	r := database.CreateClient(constants.RedisDBRateLimit)
	defer func() {
		if err := database.CloseClient(r); err != nil {
			// Log error but don't fail the main operation
			_ = err
		}
	}()

	key := "access:" + shortCode
	val, err := database.Get(r, key)
	if err != nil {
		return 0, err
	}

	// Parse the string value to int64
	var count int64
	_, err = fmt.Sscanf(val, "%d", &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

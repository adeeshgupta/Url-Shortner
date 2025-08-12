package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/adeesh/url-shortener/internal/config"
	"github.com/adeesh/url-shortener/internal/constants"
	"github.com/adeesh/url-shortener/internal/database"
	"github.com/adeesh/url-shortener/internal/utils"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// URLService handles URL shortening business logic.
type URLService struct {
	config *config.Config // Application configuration
}

// NewURLService creates a new URL service instance.
func NewURLService(cfg *config.Config) *URLService {
	return &URLService{
		config: cfg,
	}
}

// ShortenURLRequest represents the request for shortening a URL.
type ShortenURLRequest struct {
	URL         string        `json:"url"`    // The original URL to be shortened
	CustomShort string        `json:"short"`  // Optional custom short code
	Expiry      time.Duration `json:"expiry"` // Expiry time in hours
}

// ShortenURLResponse represents the response for shortening a URL.
type ShortenURLResponse struct {
	URL             string        `json:"url"`              // The original URL
	CustomShort     string        `json:"short"`            // The complete shortened URL
	Expiry          time.Duration `json:"expiry"`           // Expiry time in hours
	XRateRemaining  int           `json:"rate_limit"`       // Remaining API requests
	XRateLimitReset time.Duration `json:"rate_limit_reset"` // Time until rate limit resets
}

// ShortenURL handles the URL shortening process.
// This is the main business logic function for URL shortening.
// It performs validation, generates short codes, and persists the mapping.
// Returns a response with the shortened URL or an error.
func (s *URLService) ShortenURL(req *ShortenURLRequest) (*ShortenURLResponse, error) {
	// Validate the provided URL
	if err := s.validateURL(req.URL); err != nil {
		return nil, err
	}

	// Enforce HTTP scheme for consistency
	req.URL = utils.EnforceHTTP(req.URL)

	// Generate short code (custom or random)
	shortCode := s.generateShortCode(req.CustomShort)

	// Check if short code is available
	if err := s.checkShortCodeAvailability(shortCode); err != nil {
		return nil, err
	}

	// Set default expiry if not provided
	req.Expiry = s.setDefaultExpiry(req.Expiry)

	// Save URL mapping to database
	if err := s.saveURLMapping(shortCode, req.URL, req.Expiry); err != nil {
		return nil, err
	}

	// Build and return response
	response := s.buildResponse(req, shortCode)
	return response, nil
}

// GetOriginalURL retrieves the original URL from Redis using the short code.
func (s *URLService) GetOriginalURL(shortCode string) (string, error) {
	r := database.CreateClient(constants.RedisDBURLMappings)
	defer func() {
		if err := database.CloseClient(r); err != nil {
			// Log error but don't fail the main operation
			_ = err
		}
	}()

	value, err := database.Get(r, shortCode)
	if errors.Is(err, redis.Nil) {
		// Short code not found in database
		return "", fmt.Errorf("not found: %s", constants.ShortUrlNotFoundOnDatabase)
	} else if err != nil {
		// Database connection or other error
		return "", fmt.Errorf("database error: %s", constants.CannotConnectToTheDB)
	}
	return value, nil
}

// validateURL checks if the provided URL is valid and not the application domain (prevents infinite loops)
func (s *URLService) validateURL(url string) error {
	if !govalidator.IsURL(url) {
		return fmt.Errorf("invalid url: %s", constants.ErrorInvalidURL)
	}

	if !utils.RemoveDomainError(url) {
		return fmt.Errorf("domain error: %s", constants.ErrorInvalidURL)
	}

	return nil
}

// generateShortCode creates a short code for the URL.
func (s *URLService) generateShortCode(customShort string) string {
	if customShort == "" {
		return uuid.New().String()[:6] // Generate random 6-character code
	}
	return customShort
}

// checkShortCodeAvailability verifies if the short code is already in use.
func (s *URLService) checkShortCodeAvailability(shortCode string) error {
	r := database.CreateClient(constants.RedisDBURLMappings)
	defer func() {
		if err := database.CloseClient(r); err != nil {
			// Log error but don't fail the main operation
			_ = err
		}
	}()

	val, _ := database.Get(r, shortCode)
	if val != "" {
		return fmt.Errorf("short code in use: %s", constants.ErrorURLShortInUse)
	}

	return nil
}

// setDefaultExpiry sets default expiry if not provided.
func (s *URLService) setDefaultExpiry(expiry time.Duration) time.Duration {
	if expiry == 0 {
		return time.Duration(constants.DefaultURLExpiryHours) * time.Hour
	}
	return expiry
}

// saveURLMapping stores the URL mapping in Redis.
func (s *URLService) saveURLMapping(shortCode, originalURL string, expiry time.Duration) error {
	r := database.CreateClient(constants.RedisDBURLMappings) // Use DB 0 for URL mappings
	defer func() {
		if err := database.CloseClient(r); err != nil {
			// Log error but don't fail the main operation
			_ = err
		}
	}()

	return database.Set(r, shortCode, originalURL, expiry*3600*time.Second)
}

// buildResponse creates the response object.
func (s *URLService) buildResponse(req *ShortenURLRequest, shortCode string) *ShortenURLResponse {
	return &ShortenURLResponse{
		URL:         req.URL,
		CustomShort: s.config.Domain + "/" + shortCode,
		Expiry:      req.Expiry,
		// Rate limit fields will be populated by the handler
		XRateRemaining:  0,
		XRateLimitReset: 0,
	}
}

package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/adeesh/url-shortener/internal/constants"
	"github.com/gofiber/fiber/v2"

	"github.com/adeesh/url-shortener/internal/config"
	"github.com/adeesh/url-shortener/internal/database"
	"github.com/go-redis/redis/v8"
)

type RateLimitService struct {
	config *config.Config
}

// NewRateLimitService creates a new rate limit service instance.
func NewRateLimitService(cfg *config.Config) *RateLimitService {
	return &RateLimitService{
		config: cfg,
	}
}

// RateLimitInfo contains rate limiting information.
type RateLimitInfo struct {
	Remaining int           // Number of remaining requests allowed
	Reset     time.Duration // Time until the rate limit window resets
}

// CheckRateLimit validates if the client has exceeded rate limits.
func (s *RateLimitService) CheckRateLimit(clientIP string) error {
	r2 := database.CreateClient(constants.RedisDBRateLimit)
	defer func() {
		if err := database.CloseClient(r2); err != nil {
			// Log error but don't fail the main operation
			_ = err
		}
	}()

	val, err := database.Get(r2, clientIP)
	if errors.Is(err, redis.Nil) {
		// First request from this IP, set initial quota
		quota := strconv.Itoa(s.config.APIQuota)
		return database.Set(r2, clientIP, quota, s.config.RateLimit)
	} else if err != nil {
		return err
	}

	// Check remaining requests
	valInt, _ := strconv.Atoi(val)
	if valInt <= 0 {
		return fiber.NewError(fiber.StatusForbidden, constants.ErrorRateLimitExceeded)
	}

	return nil
}

// DecrementRateLimit decrements the rate limit counter and returns updated values.
func (s *RateLimitService) DecrementRateLimit(clientIP string) (*RateLimitInfo, error) {
	r2 := database.CreateClient(constants.RedisDBRateLimit)
	defer func() {
		if err := database.CloseClient(r2); err != nil {
			// Log error but don't fail the main operation
			_ = err
		}
	}()

	// Decrement the rate limit counter
	if err := database.Decrement(r2, clientIP); err != nil {
		return nil, err
	}

	// Get remaining requests
	val, err := database.Get(r2, clientIP)
	if err != nil {
		return nil, err
	}

	remaining, _ := strconv.Atoi(val)

	// Get time until rate limit resets
	ttl, err := database.GetTTL(r2, clientIP)
	if err != nil {
		return &RateLimitInfo{Remaining: remaining}, nil
	}

	return &RateLimitInfo{
		Remaining: remaining,
		Reset:     ttl / time.Nanosecond / time.Minute,
	}, nil
}

// GetRateLimitResetTime returns the time remaining until rate limit resets.
func (s *RateLimitService) GetRateLimitResetTime(clientIP string) (time.Duration, error) {
	r2 := database.CreateClient(constants.RedisDBRateLimit)
	defer func() {
		if err := database.CloseClient(r2); err != nil {
			// Log error but don't fail the main operation
			_ = err
		}
	}()

	ttl, err := database.GetTTL(r2, clientIP)
	if err != nil {
		return 0, err
	}

	return ttl / time.Nanosecond / time.Minute, nil
}

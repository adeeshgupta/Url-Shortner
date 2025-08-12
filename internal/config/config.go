package config

import (
	"os"
	"strconv"
	"time"

	"github.com/adeesh/url-shortener/internal/constants"
)

// Config holds all configuration values for the application.
// All configuration is loaded from environment variables with fallback defaults.
type Config struct {
	AppPort   string        // Application port (e.g., "3000")
	DBAddr    string        // Redis database address (e.g., "localhost:6379")
	DBPass    string        // Redis database password
	Domain    string        // Application domain for generating short URLs
	APIQuota  int           // Number of API requests allowed per time window
	RateLimit time.Duration // Duration of the rate limiting window
}

// Load loads configuration from environment variables with fallback defaults.
func Load() *Config {
	return &Config{
		AppPort:   getAppPort(),
		DBAddr:    getDBAddr(),
		DBPass:    getDBPass(),
		Domain:    getDomain(),
		APIQuota:  getAPIQuota(),
		RateLimit: getRateLimit(),
	}
}

// getAppPort returns the application port from environment variables.
// Defaults to "3000" if APP_PORT is not set.
func getAppPort() string {
	port := os.Getenv(constants.EnvAppPort)
	if port == "" {
		port = constants.DefaultAppPort // Default port for development
	}
	return port
}

// getDBAddr returns the Redis database address from environment variables.
// Defaults to "localhost:6379" if DB_ADDR is not set.
func getDBAddr() string {
	addr := os.Getenv(constants.EnvDBAddr)
	if addr == "" {
		addr = constants.DefaultRedisAddr // Default Redis address
	}
	return addr
}

// getDBPass returns the Redis database password from environment variables.
// Returns empty string if DB_PASS is not set.
func getDBPass() string {
	return os.Getenv(constants.EnvDBPass)
}

// getDomain returns the application domain from environment variables.
// This is used for generating complete short URLs.
// Defaults to "http://localhost:3000" if DOMAIN is not set.
func getDomain() string {
	domain := os.Getenv(constants.EnvDomain)
	if domain == "" {
		domain = constants.DefaultDomain
	}
	return domain
}

// getAPIQuota returns the API quota (requests per time window) from environment variables.
// This controls how many requests a client can make within the rate limit window.
// Defaults to 20 if API_QUOTA is not set or invalid.
func getAPIQuota() int {
	quota := os.Getenv(constants.EnvAPIQuota)
	if quota == "" {
		return constants.DefaultAPIQuota // Default quota: 10 requests per window
	}
	if quotaInt, err := strconv.Atoi(quota); err == nil {
		return quotaInt
	}
	return constants.DefaultAPIQuota
}

// getRateLimit returns the rate limit duration from environment variables.
// This controls how long the rate limiting window lasts.
// Defaults to 30 minutes if RATE_LIMIT_MINUTES is not set or invalid.
func getRateLimit() time.Duration {
	rateLimit := os.Getenv(constants.EnvRateLimitMinutes)
	if rateLimit == "" {
		return constants.DefaultRateLimitDuration // Default: 30-minute window
	}
	if rateLimitInt, err := strconv.Atoi(rateLimit); err == nil {
		return time.Duration(rateLimitInt) * time.Minute
	}
	return constants.DefaultRateLimitDuration
}

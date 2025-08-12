package constants

import "time"

// Server Configuration Constants
const (
	DefaultAppPort   = "3000"
	DefaultDomain    = "http://localhost:3000"
	DefaultRedisAddr = "localhost:6379"
)

// Rate Limiting Constants
const (
	DefaultAPIQuota          = 20
	DefaultRateLimitDuration = 30 * time.Minute
)

// URL Expiry Constants
const (
	DefaultURLExpiryHours = 24
)

// Error Messages
const (
	ErrorCannotParseJSON       = "cannot parse JSON"
	ErrorRateLimitExceeded     = "Rate limit exceeded"
	ErrorInvalidURL            = "Invalid URL"
	ErrorURLShortInUse         = "URL short already in use"
	ErrorUpdateRateLimitFailed = "Failed to update rate limit"
	ShortUrlNotFoundOnDatabase = "Short Url not found on database"
	CannotConnectToTheDB       = "Cannot connect to the DB"
)

// Redis Database Numbers
const (
	// RedisDBURLMappings is the Redis database number for URL mappings
	RedisDBURLMappings = 0
	// RedisDBRateLimit is the Redis database number for rate limiting data
	RedisDBRateLimit = 1
)

// Environment Variable Names
const (
	// EnvAppPort is the environment variable name for application port
	EnvAppPort = "APP_PORT"
	// EnvDomain is the environment variable name for application domain
	EnvDomain = "DOMAIN"
	// EnvDBAddr is the environment variable name for database address
	EnvDBAddr = "DB_ADDR"
	// EnvDBPass is the environment variable name for database password
	EnvDBPass = "DB_PASS"
	// EnvAPIQuota is the environment variable name for API quota
	EnvAPIQuota = "API_QUOTA"
	// EnvRateLimitMinutes is the environment variable name for rate limit minutes
	EnvRateLimitMinutes = "RATE_LIMIT_MINUTES"
)

const (
	Counter = "counter"
)

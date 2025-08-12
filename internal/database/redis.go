// Package database provides Redis database operations and connection management
package database

import (
	"context"
	"os"
	"time"

	"github.com/adeesh/url-shortener/internal/constants"
	"github.com/go-redis/redis/v8"
)

// Ctx is the default context used for all Redis operations.
var Ctx = context.Background()

// CreateClient creates a new Redis client for the specified database.
// Redis supports multiple databases (0-15 by default), and this function
// allows creating clients for different databases for different purposes.
//   - DB 0: URL mappings (short_code → original_url)
//   - DB 1: Analytics and rate limiting data
func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     getRedisAddr(),
		Password: getRedisPassword(),
		DB:       dbNo,
	})
	return rdb
}

// getRedisAddr returns the Redis address from environment variables.
// Defaults to "localhost:6379" if DB_ADDR is not set.
func getRedisAddr() string {
	addr := os.Getenv(constants.EnvDBAddr)
	if addr == "" {
		addr = constants.DefaultRedisAddr
	}
	return addr
}

// getRedisPassword returns the Redis password from environment variables.
// Returns empty string if DB_PASS is not set (for development without authentication).
func getRedisPassword() string {
	return os.Getenv(constants.EnvDBPass)
}

// Set stores a key-value pair in Redis with an optional expiry time.
func Set(client *redis.Client, key, value string, expiry time.Duration) error {
	return client.Set(Ctx, key, value, expiry).Err()
}

// Get retrieves a value from Redis using the provided key.
func Get(client *redis.Client, key string) (string, error) {
	return client.Get(Ctx, key).Result()
}

// Increment increments a counter in Redis.
// If the counter doesn't exist, Redis will create it starting from 0 → 1.
func Increment(client *redis.Client, key string) error {
	return client.Incr(Ctx, key).Err()
}

// Decrement decrements a counter in Redis.
func Decrement(client *redis.Client, key string) error {
	return client.Decr(Ctx, key).Err()
}

// GetTTL returns the time-to-live for a key.
func GetTTL(client *redis.Client, key string) (time.Duration, error) {
	return client.TTL(Ctx, key).Result()
}

// CloseClient safely closes a Redis client.
// It's important to close clients to prevent connection leaks.
func CloseClient(client *redis.Client) error {
	return client.Close()
}

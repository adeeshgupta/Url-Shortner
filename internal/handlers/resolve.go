package handlers

import (
	"errors"
	"net/http"

	"github.com/adeesh/url-shortener/internal/constants"
	"github.com/adeesh/url-shortener/internal/services"
	"github.com/gin-gonic/gin"
)

// analyticsService is a shared instance of the analytics service
var analyticsService = services.NewAnalyticsService()

// ResolveURL handles requests to short URLs and redirects to the original URL.
// This is the main handler for GET /:url requests.
func ResolveURL(c *gin.Context) {
	// Check rate limit BEFORE processing any request
	if err := rateLimitService.CheckRateLimit(c.ClientIP()); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": constants.ErrorRateLimitExceeded,
		})
		return
	}

	shortCode := c.Param("url")

	originalURL, err := urlService.GetOriginalURL(shortCode)
	if err != nil {
		var ginErr *gin.Error
		if errors.As(err, &ginErr) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": ginErr.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve URL",
		})
		return
	}

	// Update rate limit after successful resolution
	go func() {
		_, _ = rateLimitService.DecrementRateLimit(c.ClientIP())
		// Track total redirects
		_ = analyticsService.TrackRedirectCounter()
		// Track individual short URL access
		_ = analyticsService.TrackShortURLAccess(shortCode)
	}()

	// Redirect the user to the original URL
	// 301 = Moved Permanently (browser may cache the redirect)
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

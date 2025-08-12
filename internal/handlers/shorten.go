package handlers

import (
	"net/http"
	"time"

	"github.com/adeesh/url-shortener/internal/config"
	"github.com/adeesh/url-shortener/internal/constants"
	"github.com/adeesh/url-shortener/internal/services"
	"github.com/gin-gonic/gin"
)

// rateLimitService is a shared instance of the rate limit service
var rateLimitService = services.NewRateLimitService(config.Load())

// urlService is a shared instance of the URL service
var urlService = services.NewURLService(config.Load())

// ShortenURL handles URL shortening requests with rate limiting and validation.
// This is the main handler for POST /api/v1 requests.
func ShortenURL(c *gin.Context) {
	// Check rate limit BEFORE processing any request
	if err := rateLimitService.CheckRateLimit(c.ClientIP()); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Rate limit exceeded",
		})
		return
	}

	// Parse request body
	body, err := parseRequestBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.ErrorCannotParseJSON,
		})
		return
	}

	// Use URL service for URL shortening
	req := &services.ShortenURLRequest{
		URL:         body.URL,
		CustomShort: body.CustomShort,
		Expiry:      body.Expiry,
	}

	response, err := urlService.ShortenURL(req)
	if err != nil {
		// Handle Gin errors
		if ginErr, ok := err.(*gin.Error); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ginErr.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to shorten URL",
		})
		return
	}

	// Update rate limit and get updated values
	rateRemaining, rateReset, err := updateRateLimit(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.ErrorUpdateRateLimitFailed,
		})
		return
	}

	// Update the response with rate limit info
	response.XRateRemaining = rateRemaining
	response.XRateLimitReset = rateReset

	c.JSON(http.StatusOK, response)
}

// parseRequestBody parses and validates the incoming request body.
func parseRequestBody(c *gin.Context) (*services.ShortenURLRequest, error) {
	var body services.ShortenURLRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

// updateRateLimit decrements the rate limit counter and returns updated values.
// Uses the rate limit service for consistent rate limiting logic.
func updateRateLimit(c *gin.Context) (int, time.Duration, error) {
	rateLimitInfo, err := rateLimitService.DecrementRateLimit(c.ClientIP())
	if err != nil {
		return 0, 0, err
	}
	return rateLimitInfo.Remaining, rateLimitInfo.Reset, nil
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAnalytics returns the total redirect count and other analytics data.
// This endpoint provides access to service-wide analytics metrics.
func GetAnalytics(c *gin.Context) {
	if err := rateLimitService.CheckRateLimit(c.ClientIP()); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Rate limit exceeded",
		})
		return
	}

	count, err := analyticsService.GetRedirectCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve analytics data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_redirects": count,
		"message":         "Analytics data retrieved successfully",
	})
}

// GetShortURLAnalytics returns the access count for a specific short URL.
// This endpoint provides analytics for individual short URLs.
func GetShortURLAnalytics(c *gin.Context) {
	if err := rateLimitService.CheckRateLimit(c.ClientIP()); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Rate limit exceeded",
		})
		return
	}

	shortCode := c.Param("url")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Short code is required",
		})
		return
	}

	count, err := analyticsService.GetShortURLAccessCount(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve short URL analytics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code":   shortCode,
		"access_count": count,
		"message":      "Short URL analytics retrieved successfully",
	})
}

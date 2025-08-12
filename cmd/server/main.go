package main

import (
	"fmt"
	"log"

	"github.com/adeesh/url-shortener/internal/config"
	"github.com/adeesh/url-shortener/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// createGinApp creates and configures a new Gin application with middleware.
func createGinApp() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	return app
}

// setupRoutes configures the application routes for URL shortening and resolution.
//   - GET /:url - Resolves short URLs and redirects to original URLs
//   - POST /api/v1 - Creates shortened URLs from long URLs
//   - GET /api/v1/analytics - Returns total redirect analytics
//   - GET /api/v1/analytics/:url - Returns analytics for specific short URL
func setupRoutes(app *gin.Engine) {
	// Route for resolving short URLs (e.g., /abc123)
	app.GET("/:url", handlers.ResolveURL)

	// Route for creating shortened URLs
	app.POST("/api/v1", handlers.ShortenURL)

	// Analytics routes
	app.GET("/api/v1/analytics", handlers.GetAnalytics)
	app.GET("/api/v1/analytics/:url", handlers.GetShortURLAnalytics)
}

// startServer starts the Gin server on the configured port and starts listening for HTTP requests.
func startServer(app *gin.Engine, cfg *config.Config) error {
	port := cfg.AppPort
	log.Printf("Starting server on port %s", port)
	return app.Run(":" + port)
}

// main is the entry point of the application.
func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		err = fmt.Errorf("failed to load environment variables: %w", err)
		log.Printf("Warning: %v", err)
	}

	// Load application configuration with fallback values
	cfg := config.Load()

	// Create and configure Gin app with middleware
	app := createGinApp()

	// Setup application routes
	setupRoutes(app)

	// Start the HTTP server and listen for requests
	if err := startServer(app, cfg); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

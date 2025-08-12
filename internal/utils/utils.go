package utils

import (
	"os"
	"strings"

	"github.com/adeesh/url-shortener/internal/constants"
)

// EnforceHTTP ensures URLs have an HTTP scheme.
// If the URL doesn't start with "http", it prepends "http://".
func EnforceHTTP(url string) string {
	if url == "" {
		return "http://"
	}
	if len(url) < 4 {
		return "http://" + url
	}
	// Check for both http and https (case insensitive)
	lowerURL := strings.ToLower(url[:4])
	if lowerURL != "http" {
		return "http://" + url
	}
	return url
}

// RemoveDomainError checks if the URL is the same as the application domain.
// This prevents infinite loops from shortening the domain itself.
func RemoveDomainError(url string) bool {
	domain := os.Getenv(constants.EnvDomain)
	if domain == "" {
		return true
	}
	
	// First check exact match
	if url == domain {
		return false
	}
	
	// Clean both URLs and compare (preserving www for subdomain distinction)
	cleanUrl := cleanURLForComparison(url)
	cleanDomain := cleanURLForComparison(domain)
	
	return cleanUrl != cleanDomain
}

// cleanURL removes common URL prefixes and returns the clean domain.
func cleanURL(url string) string {
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)

	// Extract domain part (before first slash, ?, or #)
	if idx := strings.IndexAny(newURL, "/?#"); idx != -1 {
		newURL = newURL[:idx]
	}
	return newURL
}

// cleanURLForComparison removes common URL prefixes but preserves www for domain comparison.
func cleanURLForComparison(url string) string {
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	// Note: We don't remove www. to preserve subdomain distinction

	// Extract domain part (before first slash, ?, or #)
	if idx := strings.IndexAny(newURL, "/?#"); idx != -1 {
		newURL = newURL[:idx]
	}
	return newURL
}

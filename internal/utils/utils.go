package utils

import (
	"os"
	"strings"

	"github.com/adeesh/url-shortener/internal/constants"
)

// EnforceHTTP ensures URLs have an HTTP scheme.
// If the URL doesn't start with "http", it prepends "http://".
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

// RemoveDomainError checks if the URL is the same as the application domain.
// This prevents infinite loops from shortening the domain itself.
func RemoveDomainError(url string) bool {
	if url == os.Getenv(constants.EnvDomain) {
		return false
	}
	cleanUrl := cleanURL(url)
	return cleanUrl != os.Getenv(constants.EnvDomain)
}

// cleanURL removes common URL prefixes and returns the clean domain.
func cleanURL(url string) string {
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)

	// Extract domain part (before first slash)
	newURL = strings.Split(newURL, "/")[0]
	return newURL
}

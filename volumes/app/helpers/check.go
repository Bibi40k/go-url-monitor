package helpers

import (
	"net/http"
	"time"
)

// IsURLOnline checks if a URL is reachable
func IsURLOnline(url string) bool {
	client := http.Client{
		Timeout: 5 * time.Second, // Timeout for checking URL
	}
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

package helpers

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

// IsURLOnline checks if a URL is reachable
func IsURLOnline(url string) bool {
	if !isConnectionActive() {
		fmt.Println("No active internet connection detected.")
		return false
	}

	client := http.Client{
		Timeout: 10 * time.Second, // Timeout for checking URL
	}
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

// isConnectionActive pings a stable IP to check internet connectivity
func isConnectionActive() bool {
	pingTargets := []string{"8.8.8.8", "1.1.1.1"}

	for _, target := range pingTargets {
		cmd := exec.Command("ping", "-c", "1", "-W", "2", target) // Ping with 1 packet and 2-second timeout
		if err := cmd.Run(); err == nil {
			return true
		}
	}
	return false
}

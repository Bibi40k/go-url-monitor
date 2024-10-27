package main

import (
	"encoding/json"
	"fmt"
	"go-url-monitor/config"
	"go-url-monitor/helpers"
	"go-url-monitor/notifiers"
	"os"
	"time"
)

const stateFilePath = "./monitor_state.json"

var firstRun = true

func main() {
	config.LoadConfig()

	// Load previous state if available
	loadState()

	for {
		for i, urlConfig := range config.URLConfigs {
			fmt.Println("Checking for", urlConfig.URL)

			isOnline := helpers.IsURLOnline(urlConfig.URL)

			if isOnline {
				if urlConfig.Status == "offline" {
					sendNotification(urlConfig, "back online")
					config.URLConfigs[i].LastNotification = time.Now()
					config.URLConfigs[i].Status = "online"
					saveState() // Save state after every update
				}
			} else {
				if urlConfig.Status == "online" || time.Since(urlConfig.LastNotification) > config.NotificationInterval {
					sendNotification(urlConfig, "offline")
					config.URLConfigs[i].LastNotification = time.Now()
					config.URLConfigs[i].Status = "offline"
					saveState() // Save state after every update
				}
			}
		}
		fmt.Printf("Sleeping for %.0f minutes\n", config.CheckInterval.Minutes())
		time.Sleep(config.CheckInterval)
	}
}

// saveState saves the current state to a JSON file
func saveState() {
	stateData, err := json.Marshal(config.URLConfigs)
	if err != nil {
		fmt.Println("Error saving state:", err)
		return
	}

	err = os.WriteFile(stateFilePath, stateData, 0644)
	if err != nil {
		fmt.Println("Error writing state file:", err)
	}
}

// loadState loads the state from a JSON file if it exists
func loadState() {
	// Clear state if it's the app's first run
	if firstRun {
		if err := os.Remove(stateFilePath); err != nil && !os.IsNotExist(err) {
			fmt.Println("Error clearing state on first run:", err)
		} else {
			fmt.Println("State cleared on first run")
		}
		firstRun = false // Set firstRun to false after initial load
		return           // Skip loading if file was just cleared
	}

	// Load state if file exists
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return // State file does not exist, nothing to load
	}

	stateData, err := os.ReadFile(stateFilePath)
	if err != nil {
		fmt.Println("Error reading state file:", err)
		return
	}

	err = json.Unmarshal(stateData, &config.URLConfigs)
	if err != nil {
		fmt.Println("Error parsing state file:", err)
	}
}

func sendNotification(urlConfig config.URLConfig, status string) {
	// Load Bucharest timezone
	bucharestTimeZone, err := time.LoadLocation("Europe/Bucharest")
	if err != nil {
		fmt.Println("Error loading time zone:", err)
		return
	}

	// Get current time in Bucharest timezone
	currentTime := time.Now().In(bucharestTimeZone).Format("02-Jan-2006 15:04:05")

	// Create an email body template
	emailBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>Website Status Notification</h2>
			<p><strong>Status:</strong> %s</p>
			<p><strong>URL:</strong> <a href="%s">%s</a></p>
			<p><strong>Time:</strong> %s</p>
		</body>
		</html>`, status, urlConfig.URL, urlConfig.URL, currentTime)

	// Notify generic Discord channels
	for _, discord := range config.GenericDiscord {
		formattedMessage := notifiers.FormatDiscordMessage(discord.UserIDs, currentTime, status, urlConfig.URL)
		notifiers.SendDiscordMessage(discord, formattedMessage)
	}

	// Notify specific Discord channels for this URL
	for _, discord := range urlConfig.SpecificDiscord {
		formattedMessage := notifiers.FormatDiscordMessage(discord.UserIDs, currentTime, status, urlConfig.URL)
		notifiers.SendDiscordMessage(discord, formattedMessage)
	}

	// Notify via email with HTML template
	for _, email := range config.GenericEmail {
		notifiers.SendEmail(email, "URL "+status, emailBody)
	}

	for _, email := range urlConfig.SpecificEmail {
		notifiers.SendEmail(email, "Specific Monitoring: URL "+status, emailBody)
	}
}

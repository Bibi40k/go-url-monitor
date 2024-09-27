package main

import (
	"fmt"
	"go-url-monitor/config"
	"go-url-monitor/helpers"
	"go-url-monitor/notifiers"
	"time"
)

func main() {
	config.LoadConfig()

	checkInterval := 1 * time.Minute
	notificationInterval := 5 * time.Minute

	for {
		for i, urlConfig := range config.URLConfigs {
			isOnline := helpers.IsURLOnline(urlConfig.URL)

			if isOnline {
				if urlConfig.Status == "offline" {
					sendNotification(urlConfig, "back online")
					config.URLConfigs[i].LastNotification = time.Now()
					config.URLConfigs[i].Status = "online"
				}
			} else {
				if urlConfig.Status == "online" || time.Since(urlConfig.LastNotification) > notificationInterval {
					sendNotification(urlConfig, "offline")
					config.URLConfigs[i].LastNotification = time.Now()
					config.URLConfigs[i].Status = "offline"
				}
			}
		}
		time.Sleep(checkInterval)
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

	message := fmt.Sprintf("Status: %s\nURL: %s\nTime: %s", status, urlConfig.URL, currentTime)

	// Notify generic Discord channels
	for _, discord := range config.GenericDiscord {
		formattedMessage := notifiers.FormatDiscordMessage(discord.UserIDs, message, status)
		notifiers.SendDiscordMessage(discord, formattedMessage)
	}

	// Notify specific Discord channels for this URL
	for _, discord := range urlConfig.SpecificDiscord {
		formattedMessage := notifiers.FormatDiscordMessage(discord.UserIDs, message, status)
		notifiers.SendDiscordMessage(discord, formattedMessage)
	}

	// Notify via email
	for _, email := range config.GenericEmail {
		notifiers.SendEmail(email, "URL "+status, message)
	}

	for _, email := range urlConfig.SpecificEmail {
		notifiers.SendEmail(email, "Specific Monitoring: URL "+status, message)
	}
}

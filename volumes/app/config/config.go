package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// DiscordWebhook holds the configuration for Discord webhooks
type DiscordWebhook struct {
	ID      string   // Webhook ID for Discord
	Token   string   // Webhook Token for Discord
	UserIDs []string // List of user IDs to tag in the message
}

// URLConfig holds the configuration for each monitored URL
type URLConfig struct {
	URL              string           // The URL to be monitored
	SpecificEmail    []string         // List of specific email addresses for this URL
	SpecificDiscord  []DiscordWebhook // List of specific Discord webhooks for this URL
	LastNotification time.Time        // The time when the last notification was sent
	Status           string           // The current status of the URL ("online", "offline")
}

var GenericEmail []string
var GenericDiscord []DiscordWebhook
var URLConfigs []URLConfig
var DevMode bool
var CheckInterval = 1 * time.Minute
var NotificationInterval = 60 * time.Minute

// LoadConfig dynamically loads the URL configurations from environment variables
func LoadConfig() error {
	// Check if running in development mode
	DevMode = os.Getenv("DEV_MODE") == "true"
	if DevMode {
		fmt.Println("Running in Development Mode: Only generic notifications will be sent.")
	}

	// Load generic configurations
	GenericEmail = loadEnvAsSlice("GENERIC_EMAIL", ",", true)
	GenericDiscord = []DiscordWebhook{
		{
			ID:      os.Getenv("GENERIC_DISCORD_ID"),                        // Channel ID
			Token:   os.Getenv("GENERIC_DISCORD_TOKEN"),                     // Channel Token
			UserIDs: loadEnvAsSlice("GENERIC_DISCORD_USER_IDS", ",", false), // User IDs
		},
	}

	// Validate Generic Email and Discord
	if len(GenericEmail) == 0 {
		return fmt.Errorf("error: GENERIC_EMAIL must be set")
	}
	if GenericDiscord[0].ID == "" || GenericDiscord[0].Token == "" {
		return fmt.Errorf("error: GENERIC_DISCORD_ID and GENERIC_DISCORD_TOKEN must be set")
	}

	// Load total sites from environment variable
	totalSitesStr := os.Getenv("TOTAL_WATCHED_SITES")
	totalSites, err := strconv.Atoi(totalSitesStr)
	if err != nil || totalSitesStr == "" {
		return fmt.Errorf("error: TOTAL_WATCHED_SITES is not set or invalid")
	}

	// Loop through each site configuration
	for i := 1; i <= totalSites; i++ {
		indexStr := strconv.Itoa(i)
		url := os.Getenv("URL_" + indexStr)
		if url == "" {
			fmt.Printf("Warning: URL_%d is not set or is empty\n", i)
			continue // Skip if URL is not provided
		}

		// Only load specific emails and Discord if not in dev mode
		var emailList []string
		var discordWebhooks []DiscordWebhook

		if !DevMode {
			emailList = loadEnvAsSlice("URL_"+indexStr+"_SPECIFIC_EMAILS", ",", false)
			discordID := os.Getenv("URL_" + indexStr + "_SPECIFIC_DISCORD_ID")
			discordToken := os.Getenv("URL_" + indexStr + "_SPECIFIC_DISCORD_TOKEN")
			userIDList := loadEnvAsSlice("URL_"+indexStr+"_SPECIFIC_DISCORD_USER_IDS", ",", false)

			// Create DiscordWebhook object only if ID and Token are provided
			if discordID != "" && discordToken != "" {
				discordWebhooks = append(discordWebhooks, DiscordWebhook{
					ID:      discordID,
					Token:   discordToken,
					UserIDs: userIDList,
				})
			}
		}

		// Add URL configuration
		URLConfigs = append(URLConfigs, URLConfig{
			URL:             url,
			SpecificEmail:   emailList,
			SpecificDiscord: discordWebhooks,
		})
	}

	// Optionally, log loaded configuration
	logConfig()

	return nil
}

// Helper function to load environment variable as a slice
func loadEnvAsSlice(key string, sep string, required bool) []string {
	value := os.Getenv(key)
	if value == "" && required {
		fmt.Printf("Warning: Environment variable %s is not set\n", key)
	}
	if value == "" {
		return []string{} // Return empty slice if not set
	}
	return strings.Split(value, sep)
}

// Log loaded configurations for debugging
func logConfig() {
	fmt.Println("Generic configuration loaded:")
	for _, email := range GenericEmail {
		fmt.Println("Generic Email:", email)
	}
	for _, discord := range GenericDiscord {
		fmt.Println("Generic Discord Webhook ID:", discord.ID)
	}

	fmt.Println("Specific URL configurations loaded:")
	for _, config := range URLConfigs {
		fmt.Println("Monitoring URL:", config.URL)
		for _, discord := range config.SpecificDiscord {
			fmt.Printf("Specific Discord Webhook ID: %s, User IDs: %v\n", discord.ID, discord.UserIDs)
		}
	}
}

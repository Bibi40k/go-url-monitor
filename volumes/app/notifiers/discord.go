package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-url-monitor/config"
	"net/http"
)

type DiscordWebhookPayload struct {
	Content string `json:"content"`
}

func SendDiscordMessage(webhook config.DiscordWebhook, message string) {
	webhookURL := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", webhook.ID, webhook.Token)

	payload := DiscordWebhookPayload{
		Content: message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON payload:", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending Discord message:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		fmt.Printf("Failed to send Discord message. Status: %s\n", resp.Status)
	} else {
		fmt.Printf("Successfully sent message to Discord webhook ID: %s\n", webhook.ID)
	}
}

func FormatDiscordMessage(userIDs []string, message string, status string) string {
	// Tag users in Discord by their ID
	tags := ""
	for _, id := range userIDs {
		tags += fmt.Sprintf("<@%s> ", id)
	}

	// Add color using Discord's Markdown syntax
	var color string
	if status == "offline" {
		color = "```diff\n-"
	} else {
		color = "```fix\n+"
	}

	return fmt.Sprintf("%s %s%s\n```", tags, color, message)
}

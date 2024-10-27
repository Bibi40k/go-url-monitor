package notifiers

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendEmail sends an email using an SMTP server
func SendEmail(to, subject, body string) {
	// Set up authentication information.
	smtpHost := os.Getenv("SMTP_HOST") // e.g., "smtp.gmail.com"
	smtpPort := os.Getenv("SMTP_PORT") // e.g., "587"
	smtpUser := os.Getenv("SMTP_USER") // Your SMTP username
	smtpPass := os.Getenv("SMTP_PASS") // Your SMTP password

	from := smtpUser    // Use the same email as the SMTP authenticated user for "From"
	replyTo := smtpUser // Optional: Set a reply-to address

	// Create the email message with headers
	msg := []byte("To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Reply-To: " + replyTo + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" + // Set content type to HTML
		"\r\n" +
		body + "\r\n")

	// Set up the recipient and sender
	recipients := []string{to}

	// Establish connection to the SMTP server
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipients, msg)
	if err != nil {
		fmt.Printf("Error sending email to %s: %v\n", to, err)
		return
	}

	fmt.Printf("Successfully sent email to %s with subject: %s\n", to, subject)
}

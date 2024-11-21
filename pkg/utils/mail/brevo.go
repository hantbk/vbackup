package mail

import (
	"fmt"
	"log"
	"os"

	"context"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendEmail(ctx context.Context, to, subject, body string) error {

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	// Get the SMTP credentials and other settings from environment variables
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp-relay.brevo.com"
	smtpPort := 587

	adminEmailAddress := os.Getenv("ADMIN_EMAIL_ADDRESS")
	if smtpUser == "" || smtpPassword == "" || adminEmailAddress == "" {
		log.Fatal("SMTP credentials or email addresses are not set properly.")
		return fmt.Errorf("SMTP credentials or email addresses are not set")
	}

	// Create a new email message
	mailer := gomail.NewMessage()

	// Set the sender and recipient
	mailer.SetHeader("From", adminEmailAddress)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)

	// Set the email body (HTML content)
	mailer.SetBody("text/html", body)

	// Set up the SMTP client with the provided credentials
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)

	// Send the email
	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

package mail

import (
	"context"
	"testing"
)

func TestSendEmail(t *testing.T) {
	// Define test parameters
	to := "test@example.com"
	subject := "Test Email Subject"
	body := "<h1>This is a test email sent from Go using gomail</h1>"

	// Use the standard library context for `SendEmail`
	ctx := context.Background()

	// Call SendEmail
	err := SendEmail(ctx, to, subject, body)
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	// If no error, log success
	t.Log("Email sent successfully!")
}

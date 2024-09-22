package sendgrid

import (
	"fmt"
	"log"

	"github.com/cooksey14/go-recipe-blog/middleware"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// InitializeSendGridClient initializes the SendGrid client with the SendGrid API key
func InitializeSendGridClient() *sendgrid.Client {
	middleware.LoadEnvConfig()
	apiKey := middleware.GetEnv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Fatal("SENDGRID_API_KEY environment variable is not set")
	}
	client := sendgrid.NewSendClient(apiKey)
	return client
}

// SendgridSendEmail sends an email via SendGrid using the provided EmailRequest
func SendgridSendEmail(toEmail, toName string) error {
	from := mail.NewEmail("Example User", "cook.colin13@gmail.com")
	to := mail.NewEmail("", toEmail)
	subject := fmt.Sprintf("hello world")
	plainTextContent := "hello world"
	htmlContent := fmt.Sprintf("<p>%s</p>", plainTextContent)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := InitializeSendGridClient()

	response, err := client.Send(message)
	if err != nil {
		log.Printf("SendGrid Send Error: %v", err)
		return err
	}

	if response.StatusCode >= 400 {
		log.Printf("SendGrid Response Error: %s", response.Body)
		return fmt.Errorf("failed to send email: %s", response.Body)
	} else {
		log.Printf("Email sent successfully to %s with status code %d", toEmail, response.StatusCode)
		return nil
	}

}

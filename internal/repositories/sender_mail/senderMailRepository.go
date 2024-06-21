package sender_mail

import (
	"context"
	"encoding/json"
	"fmt"
	"goMailer/internal/domain"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Repository interface {
	SendEmail(ctx context.Context, params domain.SenderMailParams) error
}

type SenderMailRepository struct{}

func NewSenderMailRepository() *SenderMailRepository {
	return &SenderMailRepository{}
}

func (r *SenderMailRepository) SendEmail(ctx context.Context, params domain.SenderMailParams) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		fmt.Println("SENDGRID_API_KEY no está definida en las variables de entorno")
		return nil
	}

	from := mail.NewEmail(params.From, params.From)
	subject := params.Subject
	to := mail.NewEmail(params.From, params.To)
	plainTextContent := params.Body
	htmlContent := ""
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		fmt.Println(err)
	} else {
		if response.StatusCode >= 400 {
			var sgError domain.SendGridErrorResponse
			err := json.Unmarshal([]byte(response.Body), &sgError)
			if err != nil {
				log.Fatalf("Error parsing error response: %v", err)
			}
			for _, e := range sgError.Errors {
				fmt.Printf("Error: %s (Field: %s, Help: %s)\n", e.Message, e.Field, e.Help)
			}
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
		}
	}
	return nil
}

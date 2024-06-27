package sender_mail

import (
	"context"
	"encoding/json"
	"fmt"
	"goMailer/internal/domain"
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
		return &domain.ExternalServiceError{Err: fmt.Errorf("SENDGRID_API_KEY no está definida en las variables de entorno")}
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
		return &domain.ExternalServiceError{Err: fmt.Errorf("error sending email: %v", err)}
	}

	if response.StatusCode >= 400 && response.StatusCode < 500 {
		var sgError domain.SendGridErrorResponse
		err := json.Unmarshal([]byte(response.Body), &sgError)
		if err != nil {
			return &domain.ExternalServiceError{Err: fmt.Errorf("error parsing error response: %v", err)}
		}

		var errorMsgs []string
		for _, e := range sgError.Errors {
			errorMsgs = append(errorMsgs, fmt.Sprintf("Error: %s (Field: %s, Help: %s)", e.Message, e.Field, e.Help))
		}
		return &domain.ValidationError{Message: fmt.Sprintf("sendgrid error: %v", errorMsgs)}
	} else if response.StatusCode >= 500 {
		return &domain.ExternalServiceError{Err: fmt.Errorf("sendgrid server error: status code %d", response.StatusCode)}
	}
	return nil
}

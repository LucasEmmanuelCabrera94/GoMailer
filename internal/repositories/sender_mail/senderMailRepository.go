package sender_mail

import (
	"context"
	"encoding/json"
	"fmt"
	"goMailer/internal/domain"
	"os"
	"sync"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Repository interface {
	SendEmail(ctx context.Context, params []domain.SenderMailParams) error
}

type SenderMailRepository struct{}

func NewSenderMailRepository() *SenderMailRepository {
	return &SenderMailRepository{}
}

func (r *SenderMailRepository) SendEmail(ctx context.Context, params []domain.SenderMailParams) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		fmt.Println("SENDGRID_API_KEY no está definida en las variables de entorno")
		return &domain.ExternalServiceError{Err: fmt.Errorf("SENDGRID_API_KEY no está definida en las variables de entorno")}
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var sendErrors []error

	for _, param := range params {
		wg.Add(1)
		go func(p domain.SenderMailParams) {
			defer wg.Done()

			from := mail.NewEmail(p.From, p.From)
			subject := p.Subject
			to := mail.NewEmail(p.From, p.To)
			plainTextContent := p.Body
			htmlContent := ""
			message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
			client := sendgrid.NewSendClient(apiKey)
			response, err := client.Send(message)
			if err != nil {
				mu.Lock()
				sendErrors = append(sendErrors, fmt.Errorf("error sending email to %s: %v", p.To, err))
				mu.Unlock()
				return
			}

			if response.StatusCode >= 400 && response.StatusCode < 500 {
				var sgError domain.SendGridErrorResponse
				err := json.Unmarshal([]byte(response.Body), &sgError)
				if err != nil {
					mu.Lock()
					sendErrors = append(sendErrors, fmt.Errorf("error parsing error response for %s: %v", p.To, err))
					mu.Unlock()
					return
				}

				var errorMsgs []string
				for _, e := range sgError.Errors {
					errorMsgs = append(errorMsgs, fmt.Sprintf("Error: %s (Field: %s, Help: %s)", e.Message, e.Field, e.Help))
				}
				mu.Lock()
				sendErrors = append(sendErrors, &domain.ValidationError{Message: fmt.Sprintf("sendgrid error sending email to %s: %v", p.To, errorMsgs)})
				mu.Unlock()
			} else if response.StatusCode >= 500 {
				mu.Lock()
				sendErrors = append(sendErrors, fmt.Errorf("sendgrid server error sending email to %s: status code %d", p.To, response.StatusCode))
				mu.Unlock()
			}
		}(param)
	}

	wg.Wait()

	if len(sendErrors) > 0 {
		return &domain.ExternalServiceError{Err: fmt.Errorf("errors sending emails: %v", sendErrors)}
	}

	return nil
}

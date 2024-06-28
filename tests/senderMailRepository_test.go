package tests

import (
	"context"
	"errors"
	"goMailer/internal/domain"
	"goMailer/internal/repositories/sender_mail"
	"os"
	"testing"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

func (m *MockSendClient) Send(msg *mail.SGMailV3) (*sendgrid.Client, error) {
	return &sendgrid.Client{}, nil
}

type MockSendClientFailure struct{}

func (m *MockSendClientFailure) Send(msg *mail.SGMailV3) (*sendgrid.Client, error) {
	return nil, errors.New("sendgrid send failed")
}

func TestSendEmail_Success(t *testing.T) {
	originalAPIKey := os.Getenv("SENDGRID_API_KEY")
	defer func() {
		os.Setenv("SENDGRID_API_KEY", originalAPIKey)
	}()

	repo := &sender_mail.SenderMailRepository{}
	params := []domain.SenderMailParams{
		{
			From:    "cabreraemmanuellucas@hotmail.com",
			To:      "recipient1@hotmail.com",
			Subject: "Test Subject 1",
			Body:    "Test Body 1",
		},
	}

	err := repo.SendEmail(context.Background(), params)

	assert.NoError(t, err, "Expected no error when sending emails")
}

func TestSendEmail_Error(t *testing.T) {
	repo := &sender_mail.SenderMailRepository{}

	params := []domain.SenderMailParams{
		{
			From:    "sender@example.com",
			To:      "recipient1@example.com",
			Subject: "Test Subject 1",
			Body:    "Test Body 1",
		},
		{
			From:    "sender@example.com",
			To:      "recipient2@example.com",
			Subject: "Test Subject 2",
			Body:    "Test Body 2",
		},
	}

	err := repo.SendEmail(context.Background(), params)

	assert.Error(t, err, "Expected error when sending emails")
}

func TestSendEmail_APIKeyNotDefined(t *testing.T) {
	originalApiKey := os.Getenv("SENDGRID_API_KEY")
	defer os.Setenv("SENDGRID_API_KEY", originalApiKey)

	os.Setenv("SENDGRID_API_KEY", "")

	r := sender_mail.SenderMailRepository{}
	err := r.SendEmail(context.TODO(), []domain.SenderMailParams{})

	assert.Error(t, err)
	assert.IsType(t, &domain.ExternalServiceError{}, err)
	assert.Equal(t, "SENDGRID_API_KEY no está definida en las variables de entorno", err.(*domain.ExternalServiceError).Err.Error())
}

func TestSendEmail_ErrorSendingEmail(t *testing.T) {
	os.Setenv("SENDGRID_API_KEY", "fake-api-key")

	params := []domain.SenderMailParams{
		{From: "test@example.com", To: "recipient@example.com", Subject: "Test", Body: "Test body"},
	}

	r := sender_mail.SenderMailRepository{}
	err := r.SendEmail(context.TODO(), params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error sending email to recipient@example.com")
}

type MockSendClient struct {
	SendFunc func(*mail.SGMailV3) (*rest.Response, error)
}

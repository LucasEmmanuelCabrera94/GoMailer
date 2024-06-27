package sender_mail

import (
	"context"
	"errors"
	"goMailer/internal/domain"
	"goMailer/internal/repositories/sender_mail"
	"os"
	"testing"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

type MockSendClient struct{}

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

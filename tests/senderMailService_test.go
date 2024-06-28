package tests

import (
	"context"
	"errors"
	"goMailer/internal/domain"
	"testing"

	services "goMailer/internal/services/sender_mail"

	"github.com/stretchr/testify/assert"
)

type MockSenderMailRepository struct {
	SendEmailFunc func(ctx context.Context, params []domain.SenderMailParams) error
}

func (m *MockSenderMailRepository) SendEmail(ctx context.Context, params []domain.SenderMailParams) error {
	return m.SendEmailFunc(ctx, params)
}

func TestSenderMailService_SendEmail_ValidationError(t *testing.T) {
	mockRepo := &MockSenderMailRepository{
		SendEmailFunc: func(ctx context.Context, params []domain.SenderMailParams) error {
			return &domain.ValidationError{Message: "validation error"}
		},
	}

	service := services.NewSenderMailService(mockRepo)

	err := service.SendEmail(context.TODO(), []domain.SenderMailParams{})
	assert.Error(t, err)
	assert.IsType(t, &domain.ValidationError{}, err)
}

func TestSenderMailService_SendEmail_ExternalServiceError(t *testing.T) {
	mockRepo := &MockSenderMailRepository{
		SendEmailFunc: func(ctx context.Context, params []domain.SenderMailParams) error {
			return &domain.ExternalServiceError{Err: errors.New("external service error")}
		},
	}

	service := services.NewSenderMailService(mockRepo)

	err := service.SendEmail(context.TODO(), []domain.SenderMailParams{})
	assert.Error(t, err)
	assert.IsType(t, &domain.ExternalServiceError{}, err)
}

func TestSenderMailService_SendEmail_OtherError(t *testing.T) {
	mockRepo := &MockSenderMailRepository{
		SendEmailFunc: func(ctx context.Context, params []domain.SenderMailParams) error {
			return errors.New("some other error")
		},
	}

	service := services.NewSenderMailService(mockRepo)

	err := service.SendEmail(context.TODO(), []domain.SenderMailParams{})
	assert.Error(t, err)
	assert.IsType(t, &domain.ExternalServiceError{}, err)
	assert.Equal(t, "some other error", err.(*domain.ExternalServiceError).Err.Error())
}

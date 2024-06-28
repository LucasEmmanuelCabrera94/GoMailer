package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"goMailer/internal/api/controllers"
	"goMailer/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockSenderMailService struct {
	SendEmailFunc func(ctx context.Context, params []domain.SenderMailParams) error
}

func (m *MockSenderMailService) SendEmail(ctx context.Context, params []domain.SenderMailParams) error {
	return m.SendEmailFunc(ctx, params)
}

func TestSenderMailController_SenderMail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockSenderMailService{
		SendEmailFunc: func(ctx context.Context, params []domain.SenderMailParams) error {
			return nil
		},
	}

	controller := controllers.NewSenderMailController(mockService)
	router := gin.Default()
	router.POST("/send", controller.SenderMail)

	body, _ := json.Marshal([]domain.SenderMailParams{
		{From: "sender@example.com", To: "recipient@example.com", Subject: "Test", Body: "Test body"},
	})

	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"message":"Email sent successfully"}`, resp.Body.String())
}

func TestSenderMailController_SenderMail_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockSenderMailService{
		SendEmailFunc: func(ctx context.Context, params []domain.SenderMailParams) error {
			return &domain.ValidationError{Message: "validation error"}
		},
	}

	controller := controllers.NewSenderMailController(mockService)
	router := gin.Default()
	router.POST("/send", controller.SenderMail)

	body, _ := json.Marshal([]domain.SenderMailParams{
		{From: "invalid-email", To: "recipient@example.com", Subject: "Test", Body: "Test body"},
	})

	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid request payload")
}

func TestSenderMailController_SenderMail_ExternalServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockSenderMailService{
		SendEmailFunc: func(ctx context.Context, params []domain.SenderMailParams) error {
			return &domain.ExternalServiceError{Err: fmt.Errorf("external service error")}
		},
	}

	controller := controllers.NewSenderMailController(mockService)
	router := gin.Default()
	router.POST("/send", controller.SenderMail)

	body, _ := json.Marshal([]domain.SenderMailParams{
		{From: "sender@example.com", To: "recipient@example.com", Subject: "Test", Body: "Test body"},
	})

	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "external service error")
}

func TestSenderMailController_SenderMail_InvalidRequestPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockSenderMailService{}

	controller := controllers.NewSenderMailController(mockService)
	router := gin.Default()
	router.POST("/send", controller.SenderMail)

	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid request payload")
}

package sender_mail

import (
	"context"
	"goMailer/internal/domain"
	"log"
)

type Repository interface {
	SendEmail(ctx context.Context, params domain.SenderMailParams) error
}

type SenderMailRepository struct{}

func NewSenderMailRepository() *SenderMailRepository {
	return &SenderMailRepository{}
}

func (r *SenderMailRepository) SendEmail(ctx context.Context, params domain.SenderMailParams) error {
	log.Println("Sending email to:", params.To)

	return nil
}

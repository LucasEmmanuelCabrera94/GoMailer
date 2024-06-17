package sender_mail

import (
	"context"
	"goMailer/internal/domain"
	"goMailer/internal/repositories/sender_mail"
	"log"
)

type SenderMailService struct {
	SenderMailRepository *sender_mail.SenderMailRepository
}

func NewSenderMailService(repository *sender_mail.SenderMailRepository) *SenderMailService {
	return &SenderMailService{SenderMailRepository: repository}
}

func (svc *SenderMailService) SendEmail(ctx context.Context, params domain.SenderMailParams) error {
	log.Println("paso por Service.")
	svc.SenderMailRepository.SendEmail(ctx, params)
	return nil
}

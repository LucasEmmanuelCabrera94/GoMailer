package sender_mail

import (
	"context"
	"goMailer/internal/domain"
	"goMailer/internal/repositories/sender_mail"
)

type SenderMailService struct {
	SenderMailRepository *sender_mail.SenderMailRepository
}

func NewSenderMailService(repository *sender_mail.SenderMailRepository) *SenderMailService {
	return &SenderMailService{SenderMailRepository: repository}
}

func (svc *SenderMailService) SendEmail(ctx context.Context, params domain.SenderMailParams) error {
	err := svc.SenderMailRepository.SendEmail(ctx, params)
	if err != nil {
		switch err.(type) {
		case *domain.ValidationError:
			return err
		case *domain.ExternalServiceError:
			return err
		default:
			return &domain.ExternalServiceError{Err: err}
		}
	}
	return nil
}

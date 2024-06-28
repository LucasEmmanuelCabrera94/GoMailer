package sender_mail

import (
	"context"
	"goMailer/internal/domain"
	"goMailer/internal/repositories/sender_mail"
)

type Service interface {
	SendEmail(ctx context.Context, params []domain.SenderMailParams) error
}

type SenderMailService struct {
	SenderMailRepository sender_mail.Repository
}

func NewSenderMailService(repository sender_mail.Repository) *SenderMailService {
	return &SenderMailService{SenderMailRepository: repository}
}

func (svc *SenderMailService) SendEmail(ctx context.Context, params []domain.SenderMailParams) error {
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

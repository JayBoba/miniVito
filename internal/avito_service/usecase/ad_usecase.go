package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"mini-avito/internal/avito_service"
	"mini-avito/internal/models"
)

type adUseCase struct {
	repo avito_service.AdRepository
}

func NewAdUseCase(repo avito_service.AdRepository) avito_service.AdUseCase {
	return &adUseCase{
		repo: repo,
	}
}

func (u *adUseCase) CreateAd(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	if userID == uuid.Nil {
		return uuid.Nil, errors.New("Wrong user ID")
	}

	// TODO: логика отправки сообщения в RabbitMQ
	ad := models.Ad{
		UserID: userID,
	}

	return u.repo.CreateAd(ctx, ad)
}

func (u *adUseCase) GetAdsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Ad, error) {
	if userID == uuid.Nil {
		return nil, errors.New("Wrong user ID")
	}

	return u.repo.GetAdsByUserID(ctx, userID)
}

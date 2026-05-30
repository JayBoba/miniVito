package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"

	"mini-avito/internal/avito_service"
	"mini-avito/internal/models"
)

type adUseCase struct {
	repo      avito_service.AdRepository
	publisher AdPublisher
}

func NewAdUseCase(repo avito_service.AdRepository, publisher AdPublisher) avito_service.AdUseCase {
	return &adUseCase{
		repo:      repo,
		publisher: publisher,
	}
}

type AdPublisher interface {
	PublishAdCreated(ctx context.Context, adID string) error
}

func (u *adUseCase) CreateAd(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	if userID == uuid.Nil {
		return uuid.Nil, errors.New("invalid user ID")
	}

	ad := models.Ad{
		UserID: userID,
	}

	adID, err := u.repo.CreateAd(ctx, ad)
	if err != nil {
		return uuid.Nil, err
	}

	err = u.publisher.PublishAdCreated(ctx, adID.String())
	if err != nil {
		log.Printf("[UseCase] Error publishing to RabbitMQ: %v\n", err)
	}

	return adID, nil
}

func (u *adUseCase) GetAdsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Ad, error) {
	if userID == uuid.Nil {
		return nil, errors.New("Wrong user ID")
	}

	return u.repo.GetAdsByUserID(ctx, userID)
}

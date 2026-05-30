package avito_service

import (
	"context"

	"mini-avito/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, login, passwordHash string) (uuid.UUID, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	CreateSession(ctx context.Context, session models.Session) error
}

type UserUseCase interface {
	Register(ctx context.Context, login, password string) (uuid.UUID, error)
	Login(ctx context.Context, login, password string) (string, error)
}

type AdRepository interface {
	CreateAd(ctx context.Context, ad models.Ad) (uuid.UUID, error)
	GetAdsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Ad, error)
}

type AdUseCase interface {
	CreateAd(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetAdsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Ad, error)
}

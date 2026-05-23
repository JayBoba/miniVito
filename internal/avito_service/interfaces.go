package avito_service

import (
	"context"

	"mini-avito/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, login, passwordHash string) (int, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	CreateSession(ctx context.Context, session models.Session) error
}

type UserUseCase interface {
	Register(ctx context.Context, login, password string) (int, error)
	Login(ctx context.Context, login, password string) (string, error)
}

package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"mini-avito/internal/avito_service"
)

// userUseCase - неэкспортируемая структура!!!
type userUseCase struct {
	repo avito_service.UserRepository
}

func NewUserUseCase(repo avito_service.UserRepository) avito_service.UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (u *userUseCase) Register(ctx context.Context, login, password string) (uuid.UUID, error) {

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		return uuid.Nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userID, err := u.repo.CreateUser(ctx, login, string(hashedBytes))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user in db: %w", err)
	}
	return userID, nil
}

// Login - пока оставляем как заглушку
func (u *userUseCase) Login(ctx context.Context, login, password string) (string, error) {
	return "", fmt.Errorf("login not implemented yet")
}

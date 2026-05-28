package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"mini-avito/internal/avito_service"
	"mini-avito/internal/jwt"
	"mini-avito/internal/models"
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

func (u *userUseCase) Login(ctx context.Context, login, password string) (string, error) {
	user, err := u.repo.GetUserByLogin(ctx, login)
	if err != nil {
		log.Printf("DEBUG GetUserByLogin error: %v\n", err)
		return "", fmt.Errorf("invalid login or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("DEBUG Bcrypt error: %v\n", err)
		return "", fmt.Errorf("invalid login or password")
	}

	token, err := jwt.GenerateToken(user.ID.String())
	if err != nil {
		return "", fmt.Errorf("failed to generate jwt token: %w", err)
	}

	session := models.Session{
		SessionID: uuid.New().String(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = u.repo.CreateSession(ctx, session)
	if err != nil {
		return "", fmt.Errorf("failed to save session in db: %w", err)
	}

	return token, nil
}

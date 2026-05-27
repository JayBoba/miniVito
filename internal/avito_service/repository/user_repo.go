package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"mini-avito/internal/models"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUser(ctx context.Context, login, passwordHash string) (uuid.UUID, error) {
	var id uuid.UUID

	query := `
		INSERT INTO users (login, password) 
		VALUES ($1, $2) 
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, login, passwordHash).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *userRepo) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User

	query := `
		SELECT id, login, password, created_at 
		FROM users 
		WHERE login = $1 AND is_active = true
	`

	err := r.db.GetContext(ctx, &user, query, login)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepo) CreateSession(ctx context.Context, session models.Session) error {
	query := `
		INSERT INTO sessions (session_id, user_id, expires_at) 
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, session.ID, session.UserID, session.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

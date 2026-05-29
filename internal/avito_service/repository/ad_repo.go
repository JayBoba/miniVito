package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"mini-avito/internal/models"
)

type adRepo struct {
	db *sqlx.DB
}

func NewAdRepo(db *sqlx.DB) *adRepo {
	return &adRepo{
		db: db,
	}
}

func (r *adRepo) CreateAd(ctx context.Context, ad models.Ad) (uuid.UUID, error) {
	var id uuid.UUID

	query := `
		INSERT INTO ads (user_id, status) 
		VALUES ($1, $2) 
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, ad.UserID, ad.Status).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *adRepo) GetAdsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Ad, error) {
	var ads []models.Ad

	query := `
		SELECT id, user_id, status, created_at, completed_at 
		FROM ads 
		WHERE user_id = $1
	`

	err := r.db.SelectContext(ctx, &ads, query, userID)
	if err != nil {
		return nil, err
	}

	return ads, nil
}
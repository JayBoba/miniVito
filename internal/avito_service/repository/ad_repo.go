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
		INSERT INTO ads (user_id) 
		VALUES ($1) 
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, ad.UserID).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *adRepo) GetAdsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Ad, error) {
	var ads []models.Ad

	query := `
		SELECT id, user_id, status, created_at, completed_at, updated_at
		FROM ads 
		WHERE user_id = $1
	`

	err := r.db.SelectContext(ctx, &ads, query, userID)
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (r *adRepo) UpdateAdStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE ads SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

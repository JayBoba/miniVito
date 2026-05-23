package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SaveTestMessage(ctx context.Context, text string) error {
	query := `INSERT INTO test_messages (message) VALUES ($1)`
	_, err := r.db.ExecContext(ctx, query, text)
	return err
}
// предполагается, что в базе уже есть простая таблица test_messages с колонкой message
package models

import "time"

type User struct {
	ID        string    `db:"id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	Email     *string   `db:"email"`
	Phone     *string   `db:"phone"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Session struct {
	SessionID string    `db:"session_id"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

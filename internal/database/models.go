// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
	UserID    uuid.UUID
}

type RefreshToken struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	RevokedAt sql.NullTime
	UserID    uuid.UUID
}

type User struct {
	ID             uuid.UUID
	CreatedAt      sql.NullTime
	UpdatedAt      sql.NullTime
	Email          sql.NullString
	HashedPassword string
}

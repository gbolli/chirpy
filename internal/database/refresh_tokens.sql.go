// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: refresh_tokens.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, user_id)
VALUES (
    $1, NOW(), NOW(), $2, $3
)
RETURNING token, created_at, updated_at, expires_at, revoked_at, user_id
`

type CreateRefreshTokenParams struct {
	Token     string
	ExpiresAt time.Time
	UserID    uuid.UUID
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken, arg.Token, arg.ExpiresAt, arg.UserID)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.RevokedAt,
		&i.UserID,
	)
	return i, err
}

const getUserFromToken = `-- name: GetUserFromToken :one
SELECT user_id FROM refresh_tokens 
WHERE token = $1 
  AND expires_at > NOW() 
  AND revoked_at IS NULL
`

func (q *Queries) GetUserFromToken(ctx context.Context, token string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getUserFromToken, token)
	var user_id uuid.UUID
	err := row.Scan(&user_id)
	return user_id, err
}

const revokeToken = `-- name: RevokeToken :exec
UPDATE refresh_tokens 
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1
`

func (q *Queries) RevokeToken(ctx context.Context, token string) error {
	_, err := q.db.ExecContext(ctx, revokeToken, token)
	return err
}

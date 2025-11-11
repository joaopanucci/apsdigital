package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/infra/db"

	"github.com/google/uuid"
)

type refreshTokenRepository struct {
	db *db.PostgresDB
}

func NewRefreshTokenRepository(db *db.PostgresDB) *refreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *entities.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, $4)
	`

	token.ID = uuid.New()
	_, err := r.db.Pool.Exec(ctx, query, token.ID, token.UserID, token.Token, token.ExpiresAt)
	return err
}

func (r *refreshTokenRepository) GetByToken(ctx context.Context, token string) (*entities.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, is_revoked, created_at
		FROM refresh_tokens
		WHERE token = $1 AND is_revoked = false AND expires_at > NOW()
	`

	row := r.db.Pool.QueryRow(ctx, query, token)

	var refreshToken entities.RefreshToken
	err := row.Scan(
		&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token,
		&refreshToken.ExpiresAt, &refreshToken.IsRevoked, &refreshToken.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found or expired")
		}
		return nil, err
	}

	return &refreshToken, nil
}

func (r *refreshTokenRepository) RevokeByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET is_revoked = true WHERE user_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, userID)
	return err
}

func (r *refreshTokenRepository) RevokeToken(ctx context.Context, token string) error {
	query := `UPDATE refresh_tokens SET is_revoked = true WHERE token = $1`
	_, err := r.db.Pool.Exec(ctx, query, token)
	return err
}

func (r *refreshTokenRepository) CleanupExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.db.Pool.Exec(ctx, query)
	return err
}

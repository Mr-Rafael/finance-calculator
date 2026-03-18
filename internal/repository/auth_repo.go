package repository

import (
	"context"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthRepo struct {
	queries *db.Queries
}

func (r *AuthRepo) CreateUser(ctx context.Context, params db.CreateRefreshTokenParams) (db.RefreshToken, error) {
	return r.queries.CreateRefreshToken(ctx, params)
}

func (r *AuthRepo) GetTokenByHash(ctx context.Context, hash string) (db.RefreshToken, error) {
	return r.queries.GetTokenByHash(ctx, hash)
}

func (r *AuthRepo) RevokeTokenByUserID(ctx context.Context, id pgtype.UUID) error {
	return r.queries.RevokeTokenByUserID(ctx, id)
}

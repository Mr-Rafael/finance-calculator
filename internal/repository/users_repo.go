package repository

import (
	"context"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type UsersRepo struct {
	queries *db.Queries
}

func (r *UsersRepo) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(ctx, params)
}

func (r *UsersRepo) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}

func (r *UsersRepo) GetUserByID(ctx context.Context, id pgtype.UUID) (db.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

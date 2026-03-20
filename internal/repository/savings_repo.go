package repository

import (
	"context"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type SavingsRepo struct {
	queries *db.Queries
}

func NewSavingsRepo(queries *db.Queries) SavingsRepo {
	return SavingsRepo{queries: queries}
}

func (r *SavingsRepo) CreateSavings(ctx context.Context, params db.CreateSavingsParams) (db.Saving, error) {
	return r.queries.CreateSavings(ctx, params)
}

func (r *SavingsRepo) GetSavingsByUserID(ctx context.Context, id pgtype.UUID) ([]db.Saving, error) {
	return r.queries.GetSavingsByUserID(ctx, id)
}

func (r *SavingsRepo) CreateSavingsState(ctx context.Context, params db.CreateSavingsStateParams) (db.SavingsState, error) {
	return r.queries.CreateSavingsState(ctx, params)
}

func (r *SavingsRepo) GetSavingsStatesBySavingsID(ctx context.Context, id pgtype.UUID) ([]db.SavingsState, error) {
	return r.queries.GetSavingsStateBySavingsID(ctx, id)
}

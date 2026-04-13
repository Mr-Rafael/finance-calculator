package repository

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestSaveLoanPaymentPlan(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewAuthRepo(queries)

	test_user_id, err := uuid.Parse("af38df43-3ced-4869-9930-93a0fa0cf1e0")
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	user := pgtype.UUID{
		Bytes: test_user_id,
		Valid: true,
	}
	got, err := repo.CreateRefreshToken(ctx, user, "test_token", time.Now())
	if err != nil {
		log.Fatalf("Error saving the refresh token in database: %v", err)
	}
	want := db.RefreshToken{
		UserID: user,
	}

	if got.UserID.Bytes != want.UserID.Bytes {
		log.Fatalf("Saved (%v) and expected (%v) user IDs did not match.", got.UserID.Bytes, want.UserID.Bytes)
	}
}

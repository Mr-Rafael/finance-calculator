package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func TestCreateRefreshToken(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Printf("Error reading .env: %v", err)
	}
	dbURL := os.Getenv("POSTGRES_CONNECTION_STRING")
	ctx := context.Background()
	if dbURL == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(pool)
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

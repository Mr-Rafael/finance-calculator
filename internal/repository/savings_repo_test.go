package repository

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func TestSaveSavingsPlan(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewSavingsRepo(queries)

	test_user_id, err := uuid.Parse("af38df43-3ced-4869-9930-93a0fa0cf1e0")
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	originalData := dto.SavingsRequestParams{
		StartingCapital:     0,
		YearlyInterestRate:  "0.0",
		InterestRateType:    "APR",
		MonthlyContribution: 0,
		DurationYears:       0,
		TaxRate:             "0.0",
		YearlyInflationRate: "0.0",
		StartDate:           "1970-01-01",
	}
	params := domain.SavingsPlan{
		ID:                    uuid.Nil,
		UserID:                test_user_id,
		Name:                  "test",
		OriginalData:          domain.SavingsInput(originalData),
		StartingCapital:       decimal.Zero,
		CurrentCapital:        decimal.Zero,
		MonthlyContribution:   decimal.Zero,
		DurationMonths:        decimal.Zero,
		TaxMultiplierM:        decimal.Zero,
		InflationMultiplierY:  decimal.Zero,
		Date:                  time.Now(),
		InterestMultiplierM:   decimal.Zero,
		TotalInterestEarnings: decimal.Zero,
		RateOfReturn:          decimal.Zero,
		InflationAdjustedROR:  decimal.Zero,
	}
	status := domain.SavingsStatus{
		Date:         time.Now(),
		Interest:     0,
		Tax:          0,
		Contribution: 0,
		Increase:     0,
		Capital:      0,
	}
	params.Plan = append(params.Plan, status)

	got, err := repo.SaveSavingsPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving loan in database: %v", err)
	}

	want := db.Saving{
		UserID: pgtype.UUID{
			Bytes: test_user_id,
			Valid: true,
		},
	}

	if got.UserID.Bytes != want.UserID.Bytes {
		log.Fatalf("Saved (%v) and expected (%v) user IDs did not match.", got.UserID.Bytes, want.UserID.Bytes)
	}
}

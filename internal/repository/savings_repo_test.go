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

func TestGetSavingsPlan(t *testing.T) {
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

	plan, err := repo.SaveSavingsPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving savings plan in database: %v", err)
	}

	got, err := repo.GetSavingsPlanByID(ctx, plan.ID.Bytes, plan.UserID.Bytes)
	if err != nil {
		log.Fatalf("Error getting savings plan from database: %v", err)
	}

	want := db.Saving{
		UserID: pgtype.UUID{
			Bytes: test_user_id,
			Valid: true,
		},
	}

	if got.UserID != want.UserID.Bytes {
		log.Fatalf("The created and retrieved savings plan didn't match")
	}
}

func TestGetSavingsPlansByUser(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewSavingsRepo(queries)
	test_user_id, err := uuid.Parse("af38df43-3ced-4869-9930-93a0fa0cf1e0")
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}
	userUUID := pgtype.UUID{
		Bytes: test_user_id,
		Valid: true,
	}
	loansBefore, err := repo.queries.GetSavingsByUserID(ctx, userUUID)
	if err != nil {
		log.Fatalf("Error fetching savings plans before adding new one.")
	}
	want := len(loansBefore) + 1

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
	_, err = repo.SaveSavingsPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving loan in database: %v", err)
	}
	loansAfter, err := repo.queries.GetSavingsByUserID(ctx, userUUID)
	if err != nil {
		log.Fatalf("Error fetching loans after adding new one.")
	}
	got := len(loansAfter)

	if want != got {
		log.Fatalf("The number of loans before insert (%v) didn't match the number of loans after (%v)", want, got)
	}
}

func TestDeleteSavings(t *testing.T) {
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
	loanInfo, err := repo.SaveSavingsPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving savings plan in database: %v", err)
	}
	deleteParams := db.DeleteSavingsParams{
		ID:     loanInfo.ID,
		UserID: loanInfo.UserID,
	}
	err = repo.queries.DeleteSavings(ctx, deleteParams)
	if err != nil {
		log.Fatalf("Error deleting savings plan.")
	}

	getParams := db.GetSavingsParams{
		ID:     loanInfo.ID,
		UserID: loanInfo.UserID,
	}

	_, got := repo.queries.GetSavings(ctx, getParams)

	if got == nil {
		log.Fatalf("The savings plan was not deleted.")
	}
}

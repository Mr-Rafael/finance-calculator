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

func TestSaveLoanPaymentPlan(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewLoansRepo(queries)

	test_user_id, err := uuid.Parse("af38df43-3ced-4869-9930-93a0fa0cf1e0")
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	originalData := dto.LoanRequestParams{
		StartingPrincipal:  0,
		YearlyInterestRate: "0",
		MonthlyPayment:     0,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}
	status := domain.LoanStatus{
		Date:          time.Now(),
		Payment:       decimal.Zero,
		Interest:      decimal.Zero,
		OtherPayments: decimal.Zero,
		Paydown:       decimal.Zero,
		Principal:     decimal.Zero,
	}
	params := domain.LoanPaymentPlan{
		ID:                  uuid.Nil,
		UserID:              test_user_id,
		Name:                "test",
		OriginalData:        domain.LoansInput(originalData),
		StartingPrincipal:   decimal.Zero,
		CurrentPrincipal:    decimal.Zero,
		InterestMultiplierM: decimal.Zero,
		PaymentM:            decimal.Zero,
		EscrowM:             decimal.Zero,
		Date:                time.Now(),
		DurationMonths:      0,
		TotalExpenditure:    decimal.Zero,
		TotalPaid:           decimal.Zero,
		CostOfCreditPercent: decimal.Zero,
	}
	params.Plan = append(params.Plan, status)

	got, err := repo.SaveLoanPaymentPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving the refresh token in database: %v", err)
	}

	want := db.Loan{
		UserID: pgtype.UUID{
			Bytes: test_user_id,
			Valid: true,
		},
	}

	if got.UserID.Bytes != want.UserID.Bytes {
		log.Fatalf("Saved (%v) and expected (%v) user IDs did not match.", got.UserID.Bytes, want.UserID.Bytes)
	}
}

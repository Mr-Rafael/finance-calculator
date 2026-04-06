package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type LoansRepo struct {
	queries *db.Queries
}

func NewLoansRepo(queries *db.Queries) *LoansRepo {
	return &LoansRepo{queries: queries}
}

func (r *LoansRepo) SaveLoanPaymentPlan(ctx context.Context, plan domain.LoanPaymentPlan) (db.Loan, error) {
	loanParams, err := toLoanInsertQueryParams(plan)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Error preparing params for insert query: %v", err)
	}

	queryResult, err := r.queries.CreateLoan(ctx, loanParams)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Failed to save to database: %v", err)
	}

	for _, status := range plan.Plan {
		_, err := r.queries.CreateLoanState(ctx, toLoanStateInsertParams(status, queryResult.ID))
		if err != nil {
			return db.Loan{}, fmt.Errorf("Failed to save loan status to database: %v", err)
		}
	}
	return queryResult, nil
}

func toLoanInsertQueryParams(plan domain.LoanPaymentPlan) (db.CreateLoanParams, error) {
	startDate, err := time.Parse("2006-01-02", plan.OriginalData.StartDate)
	if err != nil {
		return db.CreateLoanParams{}, err
	}
	return db.CreateLoanParams{
		UserID: pgtype.UUID{
			Bytes: plan.UserID,
			Valid: true,
		},
		Name:               plan.Name,
		StartingPrincipal:  int32(plan.OriginalData.StartingPrincipal),
		YearlyInterestRate: plan.OriginalData.YearlyInterestRate,
		MonthlyPayment:     int32(plan.OriginalData.MonthlyPayment),
		EscrowPayment:      int32(plan.OriginalData.EscrowPayment),
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
			Valid: true,
		},
		DurationMonths:   int32(plan.DurationMonths),
		TotalExpenditure: int32(plan.TotalExpenditure.Round(0).IntPart()),
		TotalPaid:        int32(plan.TotalPaid.Round(0).IntPart()),
		CostOfCredit:     plan.CostOfCreditPercent.String(),
	}, nil
}

func toLoanStateInsertParams(status domain.LoanStatus, loanID pgtype.UUID) db.CreateLoanStateParams {
	params := db.CreateLoanStateParams{
		LoanID: loanID,
		Date: pgtype.Timestamptz{
			Time:  status.Date,
			Valid: true,
		},
		Payment:       int32(status.Payment.Round(0).IntPart()),
		Interest:      int32(status.Interest.Round(0).IntPart()),
		OtherPayments: int32(status.OtherPayments.Round(0).IntPart()),
		Paydown:       int32(status.Paydown.Round(0).IntPart()),
		Principal:     int32(status.Principal.Round(0).IntPart()),
	}
	return params
}

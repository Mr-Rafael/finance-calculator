package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type SavingsRepo struct {
	queries *db.Queries
}

func NewSavingsRepo(queries *db.Queries) *SavingsRepo {
	return &SavingsRepo{queries: queries}
}

func (r *SavingsRepo) SaveSavingsPlan(ctx context.Context, plan domain.SavingsPlan) (db.Saving, error) {
	savingsParams, err := toSavingsQueryParams(plan)
	if err != nil {
		return db.Saving{}, fmt.Errorf("Failed to save to database: %v", err)
	}

	savingsResult, err := r.queries.CreateSavings(ctx, savingsParams)
	if err != nil {
		return db.Saving{}, fmt.Errorf("Failed to save to database: %v", err)
	}

	for _, status := range plan.Plan {
		_, err := r.queries.CreateSavingsState(ctx, toSavingsStateQueryParams(status, savingsResult.ID))
		if err != nil {
			return db.Saving{}, fmt.Errorf("Failed to savings status to database: %v", err)
		}
	}
	return savingsResult, nil
}

func toSavingsQueryParams(plan domain.SavingsPlan) (db.CreateSavingsParams, error) {
	startDate, err := time.Parse("2006-01-02", plan.OriginalData.StartDate)
	if err != nil {
		return db.CreateSavingsParams{}, err
	}
	return db.CreateSavingsParams{
		UserID: pgtype.UUID{
			Bytes: plan.UserID,
			Valid: true,
		},
		Name:                plan.Name,
		StartingCapital:     int32(plan.OriginalData.StartingCapital),
		YearlyInterestRate:  plan.OriginalData.YearlyInterestRate,
		InterestRateType:    plan.OriginalData.InterestRateType,
		MonthlyContribution: int32(plan.OriginalData.MonthlyContribution),
		DurationYears:       int32(plan.OriginalData.DurationYears),
		TaxRate:             plan.OriginalData.TaxRate,
		YearlyInflationRate: pgtype.Text{
			String: plan.OriginalData.YearlyInflationRate,
			Valid:  true,
		},
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
			Valid: true,
		},
		MonthlyInterestRate:   multiplierToPercent(plan.InterestMultiplierM),
		TotalInterestEarnings: int32(plan.TotalInterestEarnings.Round(0).IntPart()),
		RateOfReturn:          plan.RateOfReturn.String(),
		InflationAdjustedRor:  plan.InflationAdjustedROR.String(),
	}, nil
}

func toSavingsStateQueryParams(status domain.SavingsStatus, savingsID pgtype.UUID) db.CreateSavingsStateParams {
	params := db.CreateSavingsStateParams{
		SavingsID: savingsID,
		Date: pgtype.Timestamptz{
			Time:  status.Date,
			Valid: true,
		},
		Interest:     int32(status.Interest),
		Tax:          int32(status.Tax),
		Contribution: int32(status.Contribution),
		Increase:     int32(status.Increase),
		Capital:      int32(status.Capital),
	}
	return params
}

func multiplierToPercent(mult decimal.Decimal) string {
	oneHundred := decimal.NewFromInt(100)
	return mult.Mul(oneHundred).String()
}

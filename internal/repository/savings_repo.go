package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/google/uuid"
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
	savingsParams, err := toSavingsInsertQueryParams(plan)
	if err != nil {
		return db.Saving{}, fmt.Errorf("Failed to save to database: %v", err)
	}

	savingsResult, err := r.queries.CreateSavings(ctx, savingsParams)
	if err != nil {
		return db.Saving{}, fmt.Errorf("Failed to save to database: %v", err)
	}

	for _, status := range plan.Plan {
		_, err := r.queries.CreateSavingsState(ctx, toSavingsStateInsertQueryParams(status, savingsResult.ID))
		if err != nil {
			return db.Saving{}, fmt.Errorf("Failed to savings status to database: %v", err)
		}
	}

	fmt.Printf("Saved this: Yearly Interest Rate: %v | Interest Rate Type: %v | Tax Rate: %v | Yearly Inflation Rate: %v\n\n",
		plan.OriginalData.YearlyInterestRate,
		plan.OriginalData.InterestRateType,
		plan.OriginalData.TaxRate,
		plan.OriginalData.YearlyInflationRate)
	return savingsResult, nil
}

func (r *SavingsRepo) GetSavingsPlansByUser(ctx context.Context, userID uuid.UUID) ([]db.GetSavingsByUserIDRow, error) {
	queryID := pgtype.UUID{
		Bytes: userID,
		Valid: true,
	}

	result, err := r.queries.GetSavingsByUserID(ctx, queryID)
	if err != nil {
		return []db.GetSavingsByUserIDRow{}, fmt.Errorf("failed to fetch user's savings plans: %v", err)
	}
	return result, nil
}

func (r *SavingsRepo) GetSavingsPlanByID(ctx context.Context, planID uuid.UUID, userID uuid.UUID) (domain.SavingsPlan, error) {
	queryID := pgtype.UUID{
		Bytes: planID,
		Valid: true,
	}

	savingsQueryResult, err := r.queries.GetSavings(ctx, toSavingsPlanSelectParams(planID, userID))
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("failed to fetch savings plan from database: %v", err)
	}
	plan, err := toSavingsPlan(savingsQueryResult)

	statesQueryResult, err := r.queries.GetSavingsStateBySavingsID(ctx, queryID)
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("failed to fetch savings plan rows from database: %v", err)
	}
	for _, state := range statesQueryResult {
		plan.Plan = append(plan.Plan, domain.SavingsStatus{
			Date:         state.Date.Time,
			Interest:     int(state.Interest),
			Tax:          int(state.Tax),
			Contribution: int(state.Contribution),
			Increase:     int(state.Increase),
			Capital:      int(state.Capital),
		})
	}

	return plan, nil
}

func toSavingsInsertQueryParams(plan domain.SavingsPlan) (db.CreateSavingsParams, error) {
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

func toSavingsStateInsertQueryParams(status domain.SavingsStatus, savingsID pgtype.UUID) db.CreateSavingsStateParams {
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

func percentToMultiplier(p string) decimal.Decimal {
	oneHundred := decimal.NewFromInt(100)
	decimalP, err := decimal.NewFromString(p)
	if err != nil {
		return decimal.Zero
	}
	return decimalP.Div(oneHundred)
}

func toSavingsPlan(queryResult db.Saving) (domain.SavingsPlan, error) {
	rateOfReturn, err := decimal.NewFromString(queryResult.RateOfReturn)
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("corrupted rate of return data for savings plan: %v", err)
	}
	inflationAdjustedReturn, err := decimal.NewFromString(queryResult.InflationAdjustedRor)
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("corrupted inflation rate of return data for savings plan: %v", err)
	}
	planData := domain.SavingsInput{
		StartingCapital:     int(queryResult.StartingCapital),
		YearlyInterestRate:  queryResult.YearlyInterestRate,
		InterestRateType:    queryResult.InterestRateType,
		MonthlyContribution: int(queryResult.MonthlyContribution),
		DurationYears:       int(queryResult.DurationYears),
		TaxRate:             queryResult.TaxRate,
		YearlyInflationRate: queryResult.YearlyInflationRate.String,
		StartDate:           queryResult.StartDate.Time.Format(time.RFC3339),
	}
	plan := domain.SavingsPlan{
		ID:                    queryResult.ID.Bytes,
		UserID:                queryResult.UserID.Bytes,
		Name:                  queryResult.Name,
		OriginalData:          planData,
		StartingCapital:       decimal.NewFromInt(int64(queryResult.StartingCapital)),
		MonthlyContribution:   decimal.NewFromInt(int64(queryResult.MonthlyContribution)),
		DurationMonths:        decimal.NewFromInt(int64(queryResult.DurationYears)).Mul(decimal.NewFromInt(12)),
		InterestMultiplierM:   percentToMultiplier(queryResult.MonthlyInterestRate),
		TotalInterestEarnings: decimal.NewFromInt(int64(queryResult.TotalInterestEarnings)),
		RateOfReturn:          rateOfReturn,
		InflationAdjustedROR:  inflationAdjustedReturn,
	}

	return plan, nil
}

func toSavingsPlanSelectParams(savingsID uuid.UUID, userID uuid.UUID) db.GetSavingsParams {
	return db.GetSavingsParams{
		ID: pgtype.UUID{
			Bytes: savingsID,
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: userID,
			Valid: true,
		},
	}
}

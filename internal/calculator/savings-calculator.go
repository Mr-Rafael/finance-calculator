package calculator

import (
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/models"
	"github.com/shopspring/decimal"
)

func CalculateSavingsPlan(info models.SavingsRequestParams) (models.SavingsPlan, error) {
	plan := models.SavingsPlan{}

	startingCapital := decimal.NewFromInt(int64(info.StartingCapital))
	currentCapital := startingCapital
	monthlyInterestRate, err := getMonthlyInterestMultiplier(info.YearlyInterestRate)
	if err != nil {
		return models.SavingsPlan{}, fmt.Errorf("failed to parse interest rate: %v", err)
	}
	durationYears := decimal.NewFromInt(int64(info.DurationYears))
	durationMonths := durationYears.Mul(decimal.NewFromInt(12))
	monthlyContribution := decimal.NewFromInt(int64(info.MonthlyContribution))
	tax := getTaxMultiplier(info.TaxRate)
	inflation := getYearlyInflationMultiplier(info.YearlyInflationRate)
	startDate, err := time.Parse("2006-01-02", info.StartDate)
	if err != nil {
		return models.SavingsPlan{}, fmt.Errorf("failed to parse start date: %v", err)
	}

	totalEarnings := decimal.NewFromInt(0)
	for i := 0; i < int(durationMonths.IntPart()); i++ {
		currentInterest := currentCapital.Mul(monthlyInterestRate)
		currentTax := currentInterest.Mul(tax)
		totalEarnings = totalEarnings.Add(currentInterest).Sub(currentTax)
		currentCapital = currentCapital.Add(currentInterest).Add(monthlyContribution).Sub(currentTax)
		currentStatus := models.SavingsStatus{
			Date:         startDate.AddDate(0, i, 0),
			Interest:     int(currentInterest.Round(0).IntPart()),
			Tax:          int(currentTax.Round(0).IntPart()),
			Contribution: int(monthlyContribution.Round(0).IntPart()),
			Increase:     int(currentInterest.Add(monthlyContribution).Sub(currentTax).Round(0).IntPart()),
			Capital:      int(currentCapital.IntPart()),
		}
		plan.Plan = append(plan.Plan, currentStatus)
	}
	rateOfReturn := currentCapital.Div(startingCapital)
	inflationAdjustedROR := rateOfReturn.Div(inflation.Pow(durationYears))

	plan.TotalInterestEarnings = int(totalEarnings.Round(0).IntPart())
	plan.RateOfReturn = getReturnPercent(rateOfReturn)
	plan.InflationAdjustedROR = getReturnPercent(inflationAdjustedROR)

	return plan, nil
}

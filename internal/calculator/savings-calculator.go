package calculator

import (
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/models"
	"github.com/shopspring/decimal"
)

type SavingsInfo struct {
	Capital             decimal.Decimal
	YearlyInterestRate  decimal.Decimal
	MonthlyContribution decimal.Decimal
	DurationYears       decimal.Decimal
	TaxRate             decimal.Decimal
	InflationRate       decimal.Decimal
	StartDate           time.Time
}

func CalculateSavingsPlan(info SavingsInfo) models.SavingsPlan {
	plan := models.SavingsPlan{}
	monthlyInterestRate := info.YearlyInterestRate.Div(decimal.NewFromInt(12))
	durationMonths := info.DurationYears.Mul(decimal.NewFromInt(12))

	currentCapital := info.Capital
	totalEarnings := decimal.NewFromInt(0)
	for i := 0; i < int(durationMonths.IntPart()); i++ {
		currentInterest := currentCapital.Mul(monthlyInterestRate)
		currentTax := currentInterest.Mul(info.TaxRate)
		totalEarnings = totalEarnings.Add(currentInterest).Sub(currentTax)
		currentCapital = currentCapital.Add(currentInterest).Add(info.MonthlyContribution).Sub(currentTax)
		currentStatus := models.SavingsStatus{
			Date:         info.StartDate.AddDate(0, i, 0),
			Interest:     int(currentInterest.IntPart()),
			Tax:          int(currentTax.IntPart()),
			Contribution: int(info.MonthlyContribution.IntPart()),
			Increase:     int(currentInterest.IntPart()) + int(info.MonthlyContribution.IntPart()),
			Capital:      int(currentCapital.IntPart()),
		}
		plan.Plan = append(plan.Plan, currentStatus)
	}
	rateOfReturn := currentCapital.Div(info.Capital).Mul(decimal.NewFromInt(100))

	plan.TotalPassiveEarnings = int(totalEarnings.IntPart())
	plan.RateOfReturn = rateOfReturn.String()
	if info.InflationRate.Cmp(decimal.NewFromInt(0)) == 1 {
		inflationMultiplier := info.InflationRate.Div(decimal.NewFromInt(100)).Add(decimal.NewFromInt(1)).Pow(info.DurationYears)
		plan.InflationAdjustedROR = rateOfReturn.Div(inflationMultiplier).String()
	}

	return plan
}

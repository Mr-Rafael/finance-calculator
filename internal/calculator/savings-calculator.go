package calculator

import (
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/models"
	"github.com/shopspring/decimal"
)

type SavingsInfo struct {
	startingCapital     decimal.Decimal
	monthlyInterestRate decimal.Decimal
	durationYears       decimal.Decimal
	durationMonths      decimal.Decimal
	monthlyContribution decimal.Decimal
	tax                 decimal.Decimal
	inflation           decimal.Decimal
	startDate           time.Time
}

func getSavingsInfoFromRequest(request models.SavingsRequestParams) (SavingsInfo, error) {
	info := SavingsInfo{}
	info.startingCapital = decimal.NewFromInt(int64(request.StartingCapital))
	monthlyInterestRate, err := getMonthlyAPYMultiplier(request.YearlyInterestRate)
	if err != nil {
		return SavingsInfo{}, fmt.Errorf("invalid interest rate: %v", request.YearlyInterestRate)
	}
	info.monthlyInterestRate = monthlyInterestRate
	info.durationYears = decimal.NewFromInt(int64(request.DurationYears))
	info.durationMonths = info.durationYears.Mul(decimal.NewFromInt(12))
	info.monthlyContribution = decimal.NewFromInt(int64(request.MonthlyContribution))
	tax, err := getTaxMultiplier(request.TaxRate)
	if err != nil {
		return SavingsInfo{}, fmt.Errorf("invalid tax rate %v", request.TaxRate)
	}
	info.tax = tax
	inflation, err := getYearlyInflationMultiplier(request.YearlyInflationRate)
	if err != nil {
		return SavingsInfo{}, fmt.Errorf("invalid inflation rate %v", request.YearlyInflationRate)
	}
	info.inflation = inflation
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return SavingsInfo{}, fmt.Errorf("invalid start date: %v", request.StartDate)
	}
	info.startDate = startDate

	return info, nil
}

func CalculateSavingsPlan(info models.SavingsRequestParams) (models.SavingsPlan, error) {
	savingsInfo, err := getSavingsInfoFromRequest(info)
	if err != nil {
		return models.SavingsPlan{}, err
	}

	plan := models.SavingsPlan{}
	currentCapital := savingsInfo.startingCapital
	totalEarnings := decimal.NewFromInt(0)

	for i := 0; i < int(savingsInfo.durationMonths.IntPart()); i++ {
		currentInterest := currentCapital.Mul(savingsInfo.monthlyInterestRate)
		currentTax := currentInterest.Mul(savingsInfo.tax)
		totalEarnings = totalEarnings.Add(currentInterest).Sub(currentTax)
		currentCapital = currentCapital.Add(currentInterest).Add(savingsInfo.monthlyContribution).Sub(currentTax)
		currentStatus := models.SavingsStatus{
			Date:         savingsInfo.startDate.AddDate(0, i, 0),
			Interest:     int(currentInterest.Round(0).IntPart()),
			Tax:          int(currentTax.Round(0).IntPart()),
			Contribution: int(savingsInfo.monthlyContribution.Round(0).IntPart()),
			Increase:     int(currentInterest.Add(savingsInfo.monthlyContribution).Sub(currentTax).Round(0).IntPart()),
			Capital:      int(currentCapital.IntPart()),
		}
		plan.Plan = append(plan.Plan, currentStatus)
	}

	rateOfReturn := totalEarnings.Add(savingsInfo.startingCapital).Div(savingsInfo.startingCapital)
	inflationAdjustedROR := rateOfReturn.Div(savingsInfo.inflation.Pow(savingsInfo.durationYears))

	plan.TotalInterestEarnings = int(totalEarnings.Round(0).IntPart())
	plan.RateOfReturn = getReturnPercent(rateOfReturn)
	plan.InflationAdjustedROR = getReturnPercent(inflationAdjustedROR)

	return plan, nil
}

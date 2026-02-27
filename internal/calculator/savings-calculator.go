package calculator

import (
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/models"
	"github.com/shopspring/decimal"
)

const minStartCapCents = "1"
const maxStartCapCents = "1000000000"
const minSavIntRate = "0.0001"
const maxSavIntRate = "1"
const minDurYears = "1"
const maxDurYears = "50"
const minMonthContrib = "0"
const maxMonthContrib = "1000000000"
const minTaxPercent = "0"
const maxTaxPercent = "100"

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
	aHundred := decimal.NewFromInt(100)

	info.startingCapital = decimal.NewFromInt(int64(request.StartingCapital))
	if !decimalIsBetween(info.startingCapital, minStartCapCents, maxStartCapCents) {
		return SavingsInfo{}, fmt.Errorf("invalid starting amount '%v'. the valid range is 0.01-1,000,000,000", info.startingCapital.Div(aHundred).Round(2))
	}

	monthlyInterestRate, err := getMonthlyInterestMultiplier(request.YearlyInterestRate, request.InterestRateType)
	if err != nil {
		return SavingsInfo{}, fmt.Errorf("invalid interest rate: %v", request.YearlyInterestRate)
	}
	if !decimalIsBetween(monthlyInterestRate, minSavIntRate, maxSavIntRate) {
		return SavingsInfo{}, fmt.Errorf("invalid interest rate '%v'. The valid range is 0.001-1", request.YearlyInterestRate)
	}
	info.monthlyInterestRate = monthlyInterestRate

	info.durationYears = decimal.NewFromInt(int64(request.DurationYears))
	info.durationMonths = info.durationYears.Mul(decimal.NewFromInt(12))
	if !decimalIsBetween(info.durationYears, minDurYears, maxDurYears) {
		return SavingsInfo{}, fmt.Errorf("invalid interest rate '%v'. The valid range is 0.001-1", request.YearlyInterestRate)
	}

	info.monthlyContribution = decimal.NewFromInt(int64(request.MonthlyContribution))
	if !decimalIsBetween(info.monthlyContribution, minMonthContrib, maxMonthContrib) {
		return SavingsInfo{}, fmt.Errorf("invalid interest rate '%v'. The valid range is 0-1,000,000,000", request.MonthlyContribution)
	}

	tax, err := getTaxMultiplier(request.TaxRate)
	if err != nil {
		return SavingsInfo{}, fmt.Errorf("invalid tax rate %v", request.TaxRate)
	}
	if !stringNumberBetween(request.TaxRate, minTaxPercent, maxTaxPercent) {
		return SavingsInfo{}, fmt.Errorf("invalid tax rate '%v'. The valid range is 0-100%%.", request.TaxRate)
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

	plan.MonthlyInterestRate = savingsInfo.monthlyInterestRate.Mul(decimal.NewFromInt(100)).Round(9).String()
	plan.TotalInterestEarnings = int(totalEarnings.Round(0).IntPart())
	plan.RateOfReturn = getReturnPercent(rateOfReturn)
	plan.InflationAdjustedROR = getReturnPercent(inflationAdjustedROR)

	return plan, nil
}

func getMonthlyInterestMultiplier(yearlyInterestRate string, interestRateType string) (decimal.Decimal, error) {
	if interestRateType == "APR" {
		return getMonthlyAPRMultiplier(yearlyInterestRate)
	}
	return getMonthlyAPYMultiplier(yearlyInterestRate)
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/Mr-Rafael/finance-calculator/internal/repository"
	"github.com/shopspring/decimal"
)

type SavingsService struct {
	repo *repository.SavingsRepo
}

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

func NewSavingsService(repo *repository.SavingsRepo) *SavingsService {
	return &SavingsService{repo: repo}
}

func (s *SavingsService) GetSavingsPlan(ctx context.Context, input domain.SavingsInput) (domain.SavingsPlan, error) {
	plan, err := initializeSavingsPlan(input)
	if err != nil {
		return domain.SavingsPlan{}, err
	}

	for i := 0; i < int(plan.DurationMonths.IntPart()); i++ {
		state := plan.PassMonth()
		state = plan.GenerateInterest(state)
		state = plan.Contribute(state)
		plan.Plan = append(plan.Plan, state)
	}
	plan.FinalCalculations()

	return plan, nil
}

func initializeSavingsPlan(input domain.SavingsInput) (domain.SavingsPlan, error) {
	plan := domain.SavingsPlan{}
	aHundred := decimal.NewFromInt(100)

	startingCapital := decimal.NewFromInt(int64(input.StartingCapital))
	if !decimalIsBetween(startingCapital, minStartCapCents, maxStartCapCents) {
		return domain.SavingsPlan{}, fmt.Errorf("invalid starting amount '%v'. the valid range is 0.01-1,000,000,000", startingCapital.Div(aHundred).Round(2))
	}
	plan.StartingCapital = startingCapital
	plan.CurrentCapital = startingCapital

	monthlyInterestRate, err := toMonthlyInterestMultiplier(input.YearlyInterestRate, input.InterestRateType)
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("invalid interest rate: %v", input.YearlyInterestRate)
	}
	if !decimalIsBetween(monthlyInterestRate, minSavIntRate, maxSavIntRate) {
		return domain.SavingsPlan{}, fmt.Errorf("invalid interest rate. The valid range is 0.001-1")
	}
	plan.InterestMultiplierM = monthlyInterestRate

	durationMonths := decimal.NewFromInt(int64(input.DurationYears)).Mul(decimal.NewFromInt(12))
	if !decimalIsBetween(durationMonths, minDurYears, maxDurYears) {
		return domain.SavingsPlan{}, fmt.Errorf("invalid plan duration. The valid range is %v-%v", minDurYears, maxDurYears)
	}
	plan.DurationMonths = durationMonths

	monthlyContribution := decimal.NewFromInt(int64(input.MonthlyContribution))
	if !decimalIsBetween(monthlyContribution, minMonthContrib, maxMonthContrib) {
		return domain.SavingsPlan{}, fmt.Errorf("invalid monthly contribution amount. The valid range is 0-1,000,000,000")
	}
	plan.MonthlyContribution = monthlyContribution

	tax, err := toTaxMultiplier(input.TaxRate)
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("invalid tax rate %v", input.TaxRate)
	}
	if !stringNumberBetween(input.TaxRate, minTaxPercent, maxTaxPercent) {
		return domain.SavingsPlan{}, fmt.Errorf("invalid tax rate '%v'. The valid range is %v-%v%%.", input.TaxRate, minTaxPercent, maxTaxPercent)
	}
	plan.TaxMultiplierM = tax

	inflation, err := toInflationMultiplier(input.YearlyInflationRate)
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("invalid inflation rate %v", input.YearlyInflationRate)
	}
	plan.InflationMultiplierY = inflation

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("invalid start date: %v", input.StartDate)
	}
	plan.Date = startDate

	return plan, nil
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/repository"
	"github.com/shopspring/decimal"
)

type SavingsService struct {
	repo *repository.SavingsRepo
}

type SavingsInput struct {
	StartingCapital     int
	YearlyInterestRate  string
	InterestRateType    string
	MonthlyContribution int
	DurationYears       int
	TaxRate             string
	YearlyInflationRate string
	StartDate           string
}

type SavingsPlan struct {
	StartingCapital       decimal.Decimal
	CurrentCapital        decimal.Decimal
	MonthlyContribution   decimal.Decimal
	DurationMonths        decimal.Decimal
	TaxMultiplierM        decimal.Decimal
	InflationMultiplierY  decimal.Decimal
	Date                  time.Time
	InterestMultiplierM   decimal.Decimal
	TotalInterestEarnings decimal.Decimal
	RateOfReturn          decimal.Decimal
	InflationAdjustedROR  decimal.Decimal
	Plan                  []SavingsStatus
}

type SavingsStatus struct {
	Date         time.Time
	Interest     int
	Tax          int
	Contribution int
	Increase     int
	Capital      int
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

func (s *SavingsService) GetSavingsPlan(ctx context.Context, input SavingsInput) (SavingsPlan, error) {
	plan, err := initializeSavingsPlan(input)
	if err != nil {
		return SavingsPlan{}, err
	}

	for i := 0; i < int(plan.DurationMonths.IntPart()); i++ {
		state := plan.passMonth()
		state = plan.generateInterest(state)
		state = plan.contribute(state)
		plan.Plan = append(plan.Plan, state)
	}
	plan.finalCalculations()

	return plan, nil
}

func initializeSavingsPlan(input SavingsInput) (SavingsPlan, error) {
	plan := SavingsPlan{}
	aHundred := decimal.NewFromInt(100)

	startingCapital := decimal.NewFromInt(int64(input.StartingCapital))
	if !decimalIsBetween(startingCapital, minStartCapCents, maxStartCapCents) {
		return SavingsPlan{}, fmt.Errorf("invalid starting amount '%v'. the valid range is 0.01-1,000,000,000", startingCapital.Div(aHundred).Round(2))
	}
	plan.StartingCapital = startingCapital
	plan.CurrentCapital = startingCapital

	monthlyInterestRate, err := toMonthlyInterestMultiplier(input.YearlyInterestRate, input.InterestRateType)
	if err != nil {
		return SavingsPlan{}, fmt.Errorf("invalid interest rate: %v", input.YearlyInterestRate)
	}
	if !decimalIsBetween(monthlyInterestRate, minSavIntRate, maxSavIntRate) {
		return SavingsPlan{}, fmt.Errorf("invalid interest rate. The valid range is 0.001-1")
	}
	plan.InterestMultiplierM = monthlyInterestRate

	durationMonths := decimal.NewFromInt(int64(input.DurationYears)).Mul(decimal.NewFromInt(12))
	if !decimalIsBetween(durationMonths, minDurYears, maxDurYears) {
		return SavingsPlan{}, fmt.Errorf("invalid plan duration. The valid range is %v-%v", minDurYears, maxDurYears)
	}
	plan.DurationMonths = durationMonths

	monthlyContribution := decimal.NewFromInt(int64(input.MonthlyContribution))
	if !decimalIsBetween(monthlyContribution, minMonthContrib, maxMonthContrib) {
		return SavingsPlan{}, fmt.Errorf("invalid monthly contribution amount. The valid range is 0-1,000,000,000")
	}
	plan.MonthlyContribution = monthlyContribution

	tax, err := toTaxMultiplier(input.TaxRate)
	if err != nil {
		return SavingsPlan{}, fmt.Errorf("invalid tax rate %v", input.TaxRate)
	}
	if !stringNumberBetween(input.TaxRate, minTaxPercent, maxTaxPercent) {
		return SavingsPlan{}, fmt.Errorf("invalid tax rate '%v'. The valid range is %v-%v%%.", input.TaxRate, minTaxPercent, maxTaxPercent)
	}
	plan.TaxMultiplierM = tax

	inflation, err := toInflationMultiplier(input.YearlyInflationRate)
	if err != nil {
		return SavingsPlan{}, fmt.Errorf("invalid inflation rate %v", input.YearlyInflationRate)
	}
	plan.InflationMultiplierY = inflation

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return SavingsPlan{}, fmt.Errorf("invalid start date: %v", input.StartDate)
	}
	plan.Date = startDate

	return plan, nil
}

func (p *SavingsPlan) passMonth() SavingsStatus {
	p.Date = p.Date.AddDate(0, 1, 0)
	return SavingsStatus{
		Date: p.Date,
	}
}

func (p *SavingsPlan) generateInterest(s SavingsStatus) SavingsStatus {
	interest := p.CurrentCapital.Mul(p.InterestMultiplierM)
	tax := interest.Mul(p.TaxMultiplierM)
	earnings := interest.Sub(tax)
	p.TotalInterestEarnings = p.TotalInterestEarnings.Add(earnings)
	p.CurrentCapital = p.CurrentCapital.Add(earnings)

	s.Interest = int(interest.Round(0).IntPart())
	s.Tax = int(tax.Round(0).IntPart())
	s.Increase = int(earnings.Round(0).IntPart())
	return s
}

func (p *SavingsPlan) contribute(s SavingsStatus) SavingsStatus {
	p.CurrentCapital = p.CurrentCapital.Add(p.MonthlyContribution)

	s.Increase = s.Increase + int(p.MonthlyContribution.Round(0).IntPart())
	s.Contribution = int(p.MonthlyContribution.Round(0).IntPart())
	s.Capital = int(p.CurrentCapital.Round(0).IntPart())
	return s
}

func (p *SavingsPlan) finalCalculations() {
	oneHundred := decimal.NewFromInt(100)
	returnRate := p.CurrentCapital.Div(p.StartingCapital)
	p.RateOfReturn = returnRate.Mul(oneHundred).Round(2)
	totalInflation := p.InflationMultiplierY.Pow(p.DurationMonths.Div(decimal.NewFromInt(12)))
	p.InflationAdjustedROR = returnRate.Div(totalInflation).Mul(oneHundred).Round(2)

}

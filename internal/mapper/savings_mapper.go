package mapper

import (
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
)

func ToSavingsResponse(plan domain.SavingsPlan) dto.SavingsResponseParams {
	response := dto.SavingsResponseParams{}

	response.MonthlyInterestRate = multiplierToPercent(plan.InterestMultiplierM)
	response.TotalInterestEarnings = int(plan.TotalInterestEarnings.Round(0).IntPart())
	response.RateOfReturn = plan.RateOfReturn.String()
	response.InflationAdjustedROR = plan.InflationAdjustedROR.String()
	for _, status := range plan.Plan {
		response.Plan = append(response.Plan, dto.SavingsStatus{
			Date:         status.Date,
			Interest:     status.Interest,
			Tax:          status.Tax,
			Contribution: status.Contribution,
			Increase:     status.Increase,
			Capital:      status.Capital,
		})
	}
	return response
}

func ToSavingsInput(input dto.SavingsRequestParams) domain.SavingsInput {
	savings := domain.SavingsInput{
		StartingCapital:     input.StartingCapital,
		YearlyInterestRate:  input.YearlyInterestRate,
		InterestRateType:    input.InterestRateType,
		MonthlyContribution: input.MonthlyContribution,
		DurationYears:       input.DurationYears,
		TaxRate:             input.TaxRate,
		YearlyInflationRate: input.YearlyInflationRate,
		StartDate:           input.StartDate,
	}

	return savings
}

package mapper

import (
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
)

func ToSavingsResponse(plan service.SavingsPlan) dto.SavingsResponseParams {
	response := dto.SavingsResponseParams{}

	response.MonthlyInterestRate = multiplierToPercent(plan.MonthlyInterestRate)
	response.TotalInterestEarnings = int(plan.TotalInterestEarnings.Round(0).IntPart())
	response.RateOfReturn = multiplierToPercent(plan.RateOfReturn)
	response.InflationAdjustedROR = multiplierToPercent(plan.InflationAdjustedROR)
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

func ToSavingsInput(input dto.SavingsRequestParams) service.SavingsInput {
	savings := service.SavingsInput{
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

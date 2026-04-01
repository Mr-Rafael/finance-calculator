package mapper

import (
	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/google/uuid"
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

func ToSavingsSaveResponse(savings db.Saving) dto.SavingsSaveResponseParams {
	return dto.SavingsSaveResponseParams{
		Name:                  savings.Name,
		ID:                    savings.ID.String(),
		MonthlyInterestRate:   savings.MonthlyInterestRate,
		TotalInterestEarnings: int(savings.TotalInterestEarnings),
		RateOfReturn:          savings.RateOfReturn,
		InflationAdjustedROR:  savings.InflationAdjustedRor,
	}
}

func ToSavingsListResponse(rows []db.GetSavingsByUserIDRow) dto.SavingsListResponseParams {
	params := dto.SavingsListResponseParams{}
	for _, row := range rows {
		newRow := dto.SavingsInfo{
			ID:              row.ID.String(),
			Name:            row.Name,
			StartingCapital: int(row.StartingCapital),
		}
		params.Plans = append(params.Plans, newRow)
	}
	return params
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

func ToSaveSavingsInput(userId uuid.UUID, input dto.SavingsSaveRequestParams) domain.SaveSavingsInput {
	savings := domain.SaveSavingsInput{
		UserID:              userId,
		PlanName:            input.Name,
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

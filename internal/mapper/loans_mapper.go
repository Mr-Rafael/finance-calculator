package mapper

import (
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
)

func ToLoanResponse(plan service.LoanPaymentPlan) dto.LoanResponseParams {
	response := dto.LoanResponseParams{}

	response.DurationMonths = int(plan.DurationMonths.Round(0).IntPart())
	response.TotalExpenditure = int(plan.TotalExpenditure.Round(0).IntPart())
	response.TotalPaid = int(plan.TotalPaid.Round(0).IntPart())
	response.CostOfCreditPercent = plan.CostOfCreditPercent.Round(2).String()
	for _, status := range plan.Plan {
		response.Plan = append(response.Plan, dto.LoanStatus{
			Date:          status.Date,
			Payment:       status.Payment,
			Interest:      status.Interest,
			OtherPayments: status.OtherPayments,
			Paydown:       status.Paydown,
			Principal:     status.Principal,
		})
	}
	return response
}

func ToLoanInput(input dto.LoanRequestParams) service.LoansInput {
	loan := service.LoansInput{
		StartingPrincipal:  input.StartingPrincipal,
		YearlyInterestRate: input.YearlyInterestRate,
		MonthlyPayment:     input.MonthlyPayment,
		EscrowPayment:      input.EscrowPayment,
		StartDate:          input.StartDate,
	}

	return loan
}

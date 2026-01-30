package calculator

import (
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/models"
	"github.com/shopspring/decimal"
)

type SavingsInfo struct {
	Capital      decimal.Decimal
	InterestRate decimal.Decimal
	Contribution decimal.Decimal
	Duration     decimal.Decimal
	StartDate    time.Time
}

func CalculateSavingsPlan(info SavingsInfo) models.SavingsPlan {
	plan := models.SavingsPlan{}

	currentCapital := info.Capital
	for i := 0; i < int(info.Duration.IntPart()); i++ {
		fmt.Printf("\nCalculating year %v:\n", i)
		fmt.Printf("Capital is %v\n", currentCapital)
		fmt.Printf("Interest is %v\n", info.InterestRate)
		fmt.Printf("Product is is %v\n", currentCapital.Mul(info.InterestRate))

		currentInterest := currentCapital.Mul(info.InterestRate)
		currentCapital = currentCapital.Add(currentInterest).Add(info.Contribution)
		currentStatus := models.SavingsStatus{
			Interest:     int(currentInterest.IntPart()),
			Contribution: int(info.Contribution.IntPart()),
			Increase:     int(currentInterest.IntPart()) + int(info.Contribution.IntPart()),
			Capital:      int(currentCapital.IntPart()),
		}
		plan.Plan = append(plan.Plan, currentStatus)
	}

	return plan
}

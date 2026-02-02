package calculator

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func getMonthlyInterestMultiplier(s string) (decimal.Decimal, error) {
	decimalInterest, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	decimalInterest = decimalInterest.Div(decimal.NewFromInt(12)).Div(decimal.NewFromInt(100))
	fmt.Printf("Interest rate parsed as: %v\n", decimalInterest)
	return decimalInterest, nil
}

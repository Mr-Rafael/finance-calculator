package calculator

import (
	"github.com/shopspring/decimal"
)

func getMonthlyInterestMultiplier(APY string) (decimal.Decimal, error) {
	decimalInterest, err := decimal.NewFromString(APY)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	decimalInterest = decimalInterest.Div(decimal.NewFromInt(100)).Add(decimal.NewFromInt(1))
	decimalInterest = decimalInterest.Pow(decimal.NewFromFloat(1.0 / 12.0)).Sub(decimal.NewFromInt(1))
	return decimalInterest, nil
}

func getTaxMultiplier(s string) decimal.Decimal {
	decimalTax, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.NewFromInt(0)
	}
	decimalTax = decimalTax.Div(decimal.NewFromInt(100))
	return decimalTax
}

func getYearlyInflationMultiplier(s string) decimal.Decimal {
	decimalInflation, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.NewFromInt(1)
	}
	decimalInflation = decimalInflation.Div(decimal.NewFromInt(100)).Add(decimal.NewFromInt(1))
	return decimalInflation
}

func getReturnPercent(rate decimal.Decimal) string {
	returnPercent := rate.Sub(decimal.NewFromInt(1)).Mul(decimal.NewFromInt(100))
	return returnPercent.Round(2).String()
}

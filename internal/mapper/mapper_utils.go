package mapper

import "github.com/shopspring/decimal"

func multiplierToPercent(mult decimal.Decimal) string {
	oneHundred := decimal.NewFromInt(100)
	two := decimal.NewFromInt(2)
	return mult.Mul(oneHundred).Div(two).Round(2).String()
}

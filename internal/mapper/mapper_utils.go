package mapper

import "github.com/shopspring/decimal"

func multiplierToPercent(mult decimal.Decimal) string {
	oneHundred := decimal.NewFromInt(100)
	return mult.Mul(oneHundred).Round(2).String()
}

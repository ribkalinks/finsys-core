package money

import (
	"github.com/shopspring/decimal"
)

//CashFlow merepresentasikan satu periode arus kas obligasi atau pinjaman.
type CashFlow struct {
	Time	decimal.Decimal // Waktu dalam tahun (Contoh: 0.5, 1.0)
	Amount	decimal.Decimal // Besaran arus kas (Kupon atau principal)
}

// CalculateMacaulayDuration menghitung Macaulay Duration dari sekumpulan arus kas.
func CalculateMacaulayDuration(cashFlows []CashFlow, yield decimal.Decimal, price decimal.Decimal) (decimal.Decimal, error) {
	if price.IsZero () {
		return decimal.Zero, ErrZeroDivision
	}

	weightedSum := decimal.Zero
	onePlusYield := decimal.NewFromInt(1).Add(yield)

	for _, cf := range cashFlows {
		//PV of cash flow = Amount / (1 + yield)^Time
		discountFactor := onePlusYield.Pow(cf.Time)
		pv := cf.Amount.Div(discountFactor)

		// Term = Time * PV
		term := cf.Time.Mul(pv)
		weightedSum = weightedSum.Add(term)
	}

	return weightedSum.Div(price), nil
}

// CalculateModifiedDuration menghitung Modified Duration dari Macaulay Duration.
func CalculateModifiedDuration(macDuration decimal.Decimal, yield decimal.Decimal) decimal.Decimal {
	onePlusYield := decimal.NewFromInt(1).Add(yield)
	return macDuration.Div(onePlusYield)
}

// CalculateConvexity menghitung Convexity dari obligasi berdasarkan arus kas.
func CalculateConvexity(cashFlows []CashFlow, yield decimal.Decimal, price decimal.Decimal) (decimal.Decimal, error) {
	if price.IsZero() {
		return decimal.Zero, ErrZeroDivision
	}

	convexitySum := decimal.Zero
	onePlusYield := decimal.NewFromInt(1).Add(yield)

	for _, cf := range cashFlows {
		// Time * (Time + 1) * PV
		discountFactor := onePlusYield.Pow(cf.Time.Add(decimal.NewFromInt(1)))
		pv := cf.Amount.Div(discountFactor)
		
		term := cf.Time.Mul(cf.Time.Add(decimal.NewFromInt(1))).Mul(pv)
		convexitySum = convexitySum.Add(term)
	}

	return convexitySum.Div(price.Mul(onePlusYield.Pow(decimal.NewFromInt(2)))), nil
}
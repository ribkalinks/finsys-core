package portfolio

import (
	"github.com/shopspring/decimal"
)

// Asset merepresentasikan sekuritas individu dengan expected return dan bobotnya
type Asset struct {
	Name           string
	ExpectedReturn decimal.Decimal
	Weight         decimal.Decimal
}

// Portfolio berisi kumpulan aset
type Portfolio struct {
	Assets []Asset
}

// CalculateExpectedReturn menghitung E(Rp) = sum(w_i * E(R_i))
func (p *Portfolio) CalculateExpectedReturn() decimal.Decimal {
	totalReturn := decimal.Zero

	for _, asset := range p.Assets {
		weightedReturn := asset.Weight.Mul(asset.ExpectedReturn)
		totalReturn = totalReturn.Add(weightedReturn)
	}

	return totalReturn
}
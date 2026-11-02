package portfolio

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestPortfolioCalculateExpectedReturn(t *testing.T) {
	tests := []struct {
		name     string
		assets   []Asset
		expected decimal.Decimal
	}{
		{
			name: "Equal weight two assets",
			assets: []Asset{
				{Name: "Asset A", ExpectedReturn: decimal.NewFromFloat(0.10), Weight: decimal.NewFromFloat(0.5)},
				{Name: "Asset B", ExpectedReturn: decimal.NewFromFloat(0.20), Weight: decimal.NewFromFloat(0.5)},
			},
			expected: decimal.NewFromFloat(0.15),
		},
		{
			name: "Unequal weights three assets",
			assets: []Asset{
				{Name: "Asset A", ExpectedReturn: decimal.NewFromFloat(0.08), Weight: decimal.NewFromFloat(0.2)},
				{Name: "Asset B", ExpectedReturn: decimal.NewFromFloat(0.12), Weight: decimal.NewFromFloat(0.5)},
				{Name: "Asset C", ExpectedReturn: decimal.NewFromFloat(0.15), Weight: decimal.NewFromFloat(0.3)},
			},
			expected: decimal.NewFromFloat(0.121), // (0.08*0.2) + (0.12*0.5) + (0.15*0.3) = 0.016 + 0.06 + 0.045 = 0.121[cite: 1]
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Portfolio{Assets: tt.assets}
			got := p.CalculateExpectedReturn()
			if !got.Equal(tt.expected) {
				t.Errorf("Expected expected return %v, got %v", tt.expected, got)
			}
		})
	}
}
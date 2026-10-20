package derivatives

import (
	"testing"
	"github.com/shopspring/decimal.v2" // sesuaikan import path decimal jika menggunakan v3 (shopspring/decimal)
)

func TestCalculateForwardPrice(t *testing.T) {
	tests := []struct {
		name      string
		spot      decimal.Decimal
		rate      decimal.Decimal
		time      decimal.Decimal
		expected  decimal.Decimal
		expectErr bool
	}{
		{
			name:     "Standard Forward Price 1 Year",
			spot:     decimal.NewFromFloat(100.0),
			rate:     decimal.NewFromFloat(0.05), // 5%
			time:     decimal.NewFromFloat(1.0),
			expected: decimal.NewFromFloat(105.0),
		},
		{
			name:     "Standard Forward Price 2 Years",
			spot:     decimal.NewFromFloat(1000.0),
			rate:     decimal.NewFromFloat(0.06), // 6%
			time:     decimal.NewFromFloat(2.0),
			expected: decimal.NewFromFloat(1123.6), // 1000 * (1.06)^2
		},
		{
			name:      "Invalid Spot Price (Zero)",
			spot:      decimal.Zero,
			rate:      decimal.NewFromFloat(0.05),
			time:      decimal.NewFromFloat(1.0),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc, err := NewForwardContract(tt.spot, tt.rate, tt.time)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			result := fc.CalculateForwardPrice()
			if !result.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCalculateForwardValue(t *testing.T) {
	currentSpot := decimal.NewFromFloat(105.0)
	initialForward := decimal.NewFromFloat(105.0)
	rate := decimal.NewFromFloat(0.05)
	timeRemaining := decimal.NewFromFloat(1.0)

	val, err := CalculateForwardValue(currentSpot, initialForward, rate, timeRemaining)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

    // currentSpot sama dengan initialForward yang telah didiskontokan pada t=1, nilainya akan positif
	if val.LessThanOrEqual(decimal.Zero) {
		t.Errorf("expected positive forward value, got %v", val)
	}
}
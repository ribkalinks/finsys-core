package portfolio

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestCalculateSampleCovariance(t *testing.T) {
	// Return historis Asset A: [10%, 20%, 30%] -> Mean = 20% (0.20)
	// Return historis Asset B: [12%, 18%, 24%] -> Mean = 18% (0.18)
	// Dev A: [-0.10, 0.00, 0.10]
	// Dev B: [-0.06, 0.00, 0.06]
	// Prod:  [0.006, 0.000, 0.006] -> Sum = 0.012
	// Cov(A, B) = 0.012 / (3 - 1) = 0.006
	returnsA := []decimal.Decimal{
		decimal.NewFromFloat(0.10),
		decimal.NewFromFloat(0.20),
		decimal.NewFromFloat(0.30),
	}
	returnsB := []decimal.Decimal{
		decimal.NewFromFloat(0.12),
		decimal.NewFromFloat(0.18),
		decimal.NewFromFloat(0.24),
	}

	expected := decimal.NewFromFloat(0.006)
	got, err := CalculateSampleCovariance(returnsA, returnsB)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !got.Equal(expected) {
		t.Errorf("Expected covariance %v, got %v", expected, got)
	}
}

func TestCalculatePortfolioVariance(t *testing.T) {
	// Portofolio 2 Aset: w = [0.6, 0.4]
	// Var(A) = 0.04, Var(B) = 0.09, Cov(A,B) = 0.01
	// σ²_p = (0.6² * 0.04) + (0.4² * 0.09) + (2 * 0.6 * 0.4 * 0.01)
	//      = 0.0144 + 0.0144 + 0.0048 = 0.0336
	weights := []decimal.Decimal{
		decimal.NewFromFloat(0.6),
		decimal.NewFromFloat(0.4),
	}
	covMatrix := [][]decimal.Decimal{
		{decimal.NewFromFloat(0.04), decimal.NewFromFloat(0.01)},
		{decimal.NewFromFloat(0.01), decimal.NewFromFloat(0.09)},
	}

	expected := decimal.NewFromFloat(0.0336)
	got, err := CalculatePortfolioVariance(weights, covMatrix)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !got.Equal(expected) {
		t.Errorf("Expected portfolio variance %v, got %v", expected, got)
	}
}
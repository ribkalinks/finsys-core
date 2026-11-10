package portfolio

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestCalculateSampleCovariance(t *testing.T) {
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

func TestCalculateBetaAndCAPM(t *testing.T) {
	assetReturns := []decimal.Decimal{
		decimal.NewFromFloat(0.12),
		decimal.NewFromFloat(0.15),
		decimal.NewFromFloat(0.08),
		decimal.NewFromFloat(0.20),
		decimal.NewFromFloat(0.14),
	}
	marketReturns := []decimal.Decimal{
		decimal.NewFromFloat(0.10),
		decimal.NewFromFloat(0.12),
		decimal.NewFromFloat(0.07),
		decimal.NewFromFloat(0.15),
		decimal.NewFromFloat(0.11),
	}

	beta, err := CalculateBeta(assetReturns, marketReturns)
	if err != nil {
		t.Fatalf("Unexpected error calculating beta: %v", err)
	}

	if beta.IsZero() {
		t.Errorf("Expected non-zero beta, got %v", beta)
	}

	riskFreeRate := decimal.NewFromFloat(0.05)
	marketExpectedReturn := decimal.NewFromFloat(0.12)

	capmReturn := CalculateCAPMExpectedReturn(riskFreeRate, marketExpectedReturn, beta)
	if capmReturn.IsZero() {
		t.Errorf("Expected non-zero CAPM return, got %v", capmReturn)
	}
}
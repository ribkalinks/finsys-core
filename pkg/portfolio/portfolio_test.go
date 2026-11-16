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

func TestRiskAdjustedPerformanceMetrics(t *testing.T) {
	portfolioReturn := decimal.NewFromFloat(0.15)      // 15%
	riskFreeRate := decimal.NewFromFloat(0.05)         // 5%
	marketExpectedReturn := decimal.NewFromFloat(0.12) // 12%
	portfolioStdDev := decimal.NewFromFloat(0.20)      // 20%
	portfolioBeta := decimal.NewFromFloat(1.2)         // 1.2

	// 1. Test Sharpe Ratio: (0.15 - 0.05) / 0.20 = 0.50
	expectedSharpe := decimal.NewFromFloat(0.5)
	gotSharpe, err := CalculateSharpeRatio(portfolioReturn, riskFreeRate, portfolioStdDev)
	if err != nil {
		t.Fatalf("Unexpected error in Sharpe Ratio: %v", err)
	}
	if !gotSharpe.Equal(expectedSharpe) {
		t.Errorf("Expected Sharpe Ratio %v, got %v", expectedSharpe, gotSharpe)
	}

	// 2. Test Treynor Ratio: (0.15 - 0.05) / 1.2 = 0.083333...
	gotTreynor, err := CalculateTreynorRatio(portfolioReturn, riskFreeRate, portfolioBeta)
	if err != nil {
		t.Fatalf("Unexpected error in Treynor Ratio: %v", err)
	}
	if gotTreynor.IsZero() {
		t.Errorf("Expected non-zero Treynor Ratio, got %v", gotTreynor)
	}

	// 3. Test Jensen's Alpha: 0.15 - [0.05 + 1.2 * (0.12 - 0.05)] = 0.15 - 0.134 = 0.016
	expectedAlpha := decimal.NewFromFloat(0.016)
	gotAlpha := CalculateJensensAlpha(portfolioReturn, riskFreeRate, marketExpectedReturn, portfolioBeta)
	if !gotAlpha.Equal(expectedAlpha) {
		t.Errorf("Expected Jensen's Alpha %v, got %v", expectedAlpha, gotAlpha)
	}
}
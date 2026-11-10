package portfolio

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// CalculateSampleCovariance menghitung kovarians antara dua himpunan data return (sampel: n-1).
func CalculateSampleCovariance(returnsA, returnsB []decimal.Decimal) (decimal.Decimal, error) {
	if len(returnsA) != len(returnsB) {
		return decimal.Zero, fmt.Errorf("panjang data return A (%d) dan B (%d) harus sama", len(returnsA), len(returnsB))
	}
	n := len(returnsA)
	if n <= 1 {
		return decimal.Zero, fmt.Errorf("jumlah data return harus lebih dari 1 untuk menghitung sampel kovarians")
	}

	sumA := decimal.Zero
	sumB := decimal.Zero
	for i := 0; i < n; i++ {
		sumA = sumA.Add(returnsA[i])
		sumB = sumB.Add(returnsB[i])
	}
	nDec := decimal.NewFromInt(int64(n))
	meanA := sumA.Div(nDec)
	meanB := sumB.Div(nDec)

	sumProduct := decimal.Zero
	for i := 0; i < n; i++ {
		devA := returnsA[i].Sub(meanA)
		devB := returnsB[i].Sub(meanB)
		sumProduct = sumProduct.Add(devA.Mul(devB))
	}

	nMinusOne := decimal.NewFromInt(int64(n - 1))
	cov := sumProduct.Div(nMinusOne)

	return cov, nil
}

// CalculatePortfolioVariance menghitung total varians portofolio berdasarkan matriks kovarians antar aset:
// σ²_p = Σ Σ w_i * w_j * Cov(i, j)
func CalculatePortfolioVariance(weights []decimal.Decimal, covMatrix [][]decimal.Decimal) (decimal.Decimal, error) {
	n := len(weights)
	if n == 0 {
		return decimal.Zero, fmt.Errorf("weights tidak boleh kosong")
	}
	if len(covMatrix) != n {
		return decimal.Zero, fmt.Errorf("ukuran matriks kovarians (%d) harus sesuai dengan jumlah weights (%d)", len(covMatrix), n)
	}

	totalVariance := decimal.Zero
	for i := 0; i < n; i++ {
		if len(covMatrix[i]) != n {
			return decimal.Zero, fmt.Errorf("matriks kovarians baris ke-%d tidak simetris (ukuran %d vs %d)", i, len(covMatrix[i]), n)
		}
		for j := 0; j < n; j++ {
			term := weights[i].Mul(weights[j]).Mul(covMatrix[i][j])
			totalVariance = totalVariance.Add(term)
		}
	}

	return totalVariance, nil
}

// CalculateBeta menghitung nilai Beta (β) aset terhadap pasar:
// β = Cov(R_i, R_m) / Var(R_m)
func CalculateBeta(assetReturns, marketReturns []decimal.Decimal) (decimal.Decimal, error) {
	cov, err := CalculateSampleCovariance(assetReturns, marketReturns)
	if err != nil {
		return decimal.Zero, fmt.Errorf("gagal menghitung kovarians aset dan pasar: %w", err)
	}

	marketVariance, err := CalculateSampleCovariance(marketReturns, marketReturns)
	if err != nil {
		return decimal.Zero, fmt.Errorf("gagal menghitung varians pasar: %w", err)
	}

	if marketVariance.IsZero() {
		return decimal.Zero, fmt.Errorf("varians pasar tidak boleh nol (tidak ada volatilitas pasar)")
	}

	beta := cov.Div(marketVariance)
	return beta, nil
}

// CalculateCAPMExpectedReturn menghitung expected return menggunakan Capital Asset Pricing Model (CAPM):
// E(R_i) = R_f + β * (E(R_m) - R_f)
func CalculateCAPMExpectedReturn(riskFreeRate, marketExpectedReturn, beta decimal.Decimal) decimal.Decimal {
	marketRiskPremium := marketExpectedReturn.Sub(riskFreeRate)
	systematicRiskReturn := beta.Mul(marketRiskPremium)
	expectedReturn := riskFreeRate.Add(systematicRiskReturn)
	return expectedReturn
}
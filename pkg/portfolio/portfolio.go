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

	// 1. Hitung Mean A dan Mean B
	sumA := decimal.Zero
	sumB := decimal.Zero
	for i := 0; i < n; i++ {
		sumA = sumA.Add(returnsA[i])
		sumB = sumB.Add(returnsB[i])
	}
	nDec := decimal.NewFromInt(int64(n))
	meanA := sumA.Div(nDec)
	meanB := sumB.Div(nDec)

	// 2. Hitung Sum of Products of Deviations: Σ(A_i - meanA)(B_i - meanB)
	sumProduct := decimal.Zero
	for i := 0; i < n; i++ {
		devA := returnsA[i].Sub(meanA)
		devB := returnsB[i].Sub(meanB)
		sumProduct = sumProduct.Add(devA.Mul(devB))
	}

	// 3. Bagi dengan (n - 1)
	nMinusOne := decimal.NewFromInt(int64(n - 1))
	cov := sumProduct.Div(nMinusOne)

	return cov, nil
}

// CalculateVariance menghitung total varians portofolio berdasarkan matriks kovarians antar aset:
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
			// w_i * w_j * Cov(i, j)
			term := weights[i].Mul(weights[j]).Mul(covMatrix[i][j])
			totalVariance = totalVariance.Add(term)
		}
	}

	return totalVariance, nil
}
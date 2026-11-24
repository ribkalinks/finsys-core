package accounting

import (
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

// TrialBalanceRow merepresentasikan satu baris akun pada Neraca Saldo
type TrialBalanceRow struct {
	AccountID   string
	AccountName string
	AccountType AccountType
	Debit       decimal.Decimal
	Credit      decimal.Decimal
}

// TrialBalance merepresentasikan rekapitulasi saldo seluruh akun pada waktu tertentu
type TrialBalance struct {
	AsOf        time.Time
	Rows        []TrialBalanceRow
	TotalDebit  decimal.Decimal
	TotalCredit decimal.Decimal
}

// IsBalanced memeriksa apakah total sisi Debit sama persis dengan total sisi Kredit
func (tb *TrialBalance) IsBalanced() bool {
	return tb.TotalDebit.Equal(tb.TotalCredit)
}

// GenerateTrialBalance menyusun TrialBalance dari seluruh akun dan saldo pada Ledger
func GenerateTrialBalance(ledger *Ledger, asOf time.Time) (*TrialBalance, error) {
	if ledger == nil {
		return nil, fmt.Errorf("ledger tidak boleh nil")
	}

	accounts := ledger.GetAllAccounts()
	// Urutkan berdasarkan AccountID agar hasil neraca saldo deterministik
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].ID < accounts[j].ID
	})

	rows := make([]TrialBalanceRow, 0, len(accounts))
	totalDebit := decimal.Zero
	totalCredit := decimal.Zero

	for _, acc := range accounts {
		balance, err := ledger.GetAccountBalance(acc.ID)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil saldo akun %s: %w", acc.ID, err)
		}

		debit := decimal.Zero
		credit := decimal.Zero

		// GetAccountBalance mengembalikan nilai positif jika sesuai NormalBalance,
		// dan negatif jika saldonya berlawanan dengan NormalBalance.
		switch acc.NormalBalance {
		case DebitBalance:
			if balance.GreaterThanOrEqual(decimal.Zero) {
				debit = balance
			} else {
				credit = balance.Abs()
			}
		case CreditBalance:
			if balance.GreaterThanOrEqual(decimal.Zero) {
				credit = balance
			} else {
				debit = balance.Abs()
			}
		}

		totalDebit = totalDebit.Add(debit)
		totalCredit = totalCredit.Add(credit)

		rows = append(rows, TrialBalanceRow{
			AccountID:   acc.ID,
			AccountName: acc.Name,
			AccountType: acc.Type,
			Debit:       debit,
			Credit:      credit,
		})
	}

	return &TrialBalance{
		AsOf:        asOf,
		Rows:        rows,
		TotalDebit:  totalDebit,
		TotalCredit: totalCredit,
	}, nil
}
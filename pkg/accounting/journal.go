package accounting

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type EntryType string

const (
	Debit  EntryType = "DEBIT"
	Credit EntryType = "CREDIT"
)

type JournalLine struct {
	AccountID string
	Type      EntryType
	Amount    decimal.Decimal
}

type JournalEntry struct {
	ID          string
	Date        time.Time
	Description string
	Lines       []JournalLine
}

// NewJournalEntry membuat entri jurnal baru sekaligus memvalidasi keseimbangan Debit == Kredit
func NewJournalEntry(id string, date time.Time, description string, lines []JournalLine) (*JournalEntry, error) {
	if len(lines) < 2 {
		return nil, fmt.Errorf("entri jurnal minimal harus memiliki 2 baris (double-entry)")
	}

	totalDebit := decimal.Zero
	totalCredit := decimal.Zero

	for _, line := range lines {
		if line.Amount.LessThanOrEqual(decimal.Zero) {
			return nil, fmt.Errorf("nominal transaksi pada akun %s harus lebih besar dari nol", line.AccountID)
		}
		switch line.Type {
		case Debit:
			totalDebit = totalDebit.Add(line.Amount)
		case Credit:
			totalCredit = totalCredit.Add(line.Amount)
		default:
			return nil, fmt.Errorf("tipe entri tidak dikenal: %s", line.Type)
		}
	}

	if !totalDebit.Equal(totalCredit) {
		return nil, fmt.Errorf("jurnal tidak seimbang (tidak balance): total debit (%s) != total kredit (%s)", totalDebit.String(), totalCredit.String())
	}

	return &JournalEntry{
		ID:          id,
		Date:        date,
		Description: description,
		Lines:       lines,
	}, nil
}
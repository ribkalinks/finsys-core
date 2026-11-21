package accounting

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// LedgerEntry merepresentasikan satu baris mutasi yang tercatat pada akun tertentu di Buku Besar
type LedgerEntry struct {
	JournalID   string
	Date        time.Time
	Description string
	Type        EntryType
	Amount      decimal.Decimal
}

// Ledger merepresentasikan Buku Besar (General Ledger) yang mengelola Chart of Accounts dan riwayat posting
type Ledger struct {
	accounts       map[string]*Account
	entries        map[string][]LedgerEntry
	postedJournals map[string]*JournalEntry
}

// NewLedger menginisialisasi instance Buku Besar kosong
func NewLedger() *Ledger {
	return &Ledger{
		accounts:       make(map[string]*Account),
		entries:        make(map[string][]LedgerEntry),
		postedJournals: make(map[string]*JournalEntry),
	}
}

// AddAccount mendaftarkan akun baru ke dalam Chart of Accounts pada Buku Besar
func (l *Ledger) AddAccount(acc *Account) error {
	if acc == nil {
		return fmt.Errorf("akun tidak boleh nil")
	}
	if _, exists := l.accounts[acc.ID]; exists {
		return fmt.Errorf("akun dengan ID %s sudah terdaftar di buku besar", acc.ID)
	}
	l.accounts[acc.ID] = acc
	l.entries[acc.ID] = []LedgerEntry{}
	return nil
}

// GetAccount mengambil informasi akun berdasarkan ID
func (l *Ledger) GetAccount(accountID string) (*Account, error) {
	acc, exists := l.accounts[accountID]
	if !exists {
		return nil, fmt.Errorf("akun dengan ID %s tidak ditemukan", accountID)
	}
	return acc, nil
}

// PostJournal memvalidasi dan memposting JournalEntry ke masing-masing akun di Buku Besar
func (l *Ledger) PostJournal(je *JournalEntry) error {
	if je == nil {
		return fmt.Errorf("entri jurnal tidak boleh nil")
	}
	if _, exists := l.postedJournals[je.ID]; exists {
		return fmt.Errorf("jurnal dengan ID %s sudah pernah diposting", je.ID)
	}

	// Validasi terlebih dahulu bahwa semua akun dalam baris jurnal terdaftar (atomicity check)
	for _, line := range je.Lines {
		if _, exists := l.accounts[line.AccountID]; !exists {
			return fmt.Errorf("gagal posting jurnal %s: akun %s tidak terdaftar di buku besar", je.ID, line.AccountID)
		}
	}

	// Catat mutasi ke masing-masing akun
	for _, line := range je.Lines {
		l.entries[line.AccountID] = append(l.entries[line.AccountID], LedgerEntry{
			JournalID:   je.ID,
			Date:        je.Date,
			Description: je.Description,
			Type:        line.Type,
			Amount:      line.Amount,
		})
	}

	l.postedJournals[je.ID] = je
	return nil
}

// GetAccountEntries mengembalikan seluruh riwayat mutasi (Debit/Credit) untuk suatu akun
func (l *Ledger) GetAccountEntries(accountID string) ([]LedgerEntry, error) {
	if _, exists := l.accounts[accountID]; !exists {
		return nil, fmt.Errorf("akun dengan ID %s tidak ditemukan", accountID)
	}
	return l.entries[accountID], nil
}

// GetAccountBalance menghitung saldo bersih suatu akun berdasarkan aturan Normal Balance
func (l *Ledger) GetAccountBalance(accountID string) (decimal.Decimal, error) {
	acc, exists := l.accounts[accountID]
	if !exists {
		return decimal.Zero, fmt.Errorf("akun dengan ID %s tidak ditemukan", accountID)
	}

	totalDebit := decimal.Zero
	totalCredit := decimal.Zero

	for _, entry := range l.entries[accountID] {
		switch entry.Type {
		case Debit:
			totalDebit = totalDebit.Add(entry.Amount)
		case Credit:
			totalCredit = totalCredit.Add(entry.Amount)
		}
	}

	if acc.NormalBalance == DebitBalance {
		return totalDebit.Sub(totalCredit), nil
	}
	return totalCredit.Sub(totalDebit), nil
}

// GetAllAccounts mengembalikan seluruh daftar akun yang terdaftar di Buku Besar
func (l *Ledger) GetAllAccounts() []*Account {
	result := make([]*Account, 0, len(l.accounts))
	for _, acc := range l.accounts {
		result = append(result, acc)
	}
	return result
}
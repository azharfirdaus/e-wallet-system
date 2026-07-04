package model

import "time"

type LedgerOperation string

const (
	LedgerOperationTopup   LedgerOperation = "TOPUP"
	LedgerOperationPayment LedgerOperation = "PAYMENT"
	LedgerOperationSend    LedgerOperation = "SEND"
	LedgerOperationReceive LedgerOperation = "RECEIVE"
)

type LedgerReference string

const (
	LedgerReferenceTopup    LedgerReference = "TOPUP"
	LedgerReferencePayment  LedgerReference = "PAYMENT"
	LedgerReferenceTransfer LedgerReference = "TRANSFER"
	LedgerReferenceReceive  LedgerReference = "RECEIVE"
)

type Ledger struct {
	ID               int64           `db:"id"`
	WalletCurrencyID int64           `db:"wallet_currency_id"`
	Debit            int64           `db:"debit"`
	Credit           int64           `db:"credit"`
	Reference        LedgerReference `db:"reference"`
	CreatedAt        time.Time       `db:"created_at"`
}

type LedgerSum struct {
	WalletCurrencyID int64 `db:"wallet_currency_id"`
	Debit            int64 `db:"debit"`
	Credit           int64 `db:"credit"`
}

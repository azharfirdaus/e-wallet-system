package model

import "time"

type LedgerOperation string

const (
	LedgerOperationTopup   LedgerOperation = "TOPUP"
	LedgerOperationPayment LedgerOperation = "PAYMENT"
	LedgerOperationSend    LedgerOperation = "SEND"
	LedgerOperationReceive LedgerOperation = "RECEIVE"
)

type Ledger struct {
	ID        int64     `db:"id"`
	Debit     int64     `db:"debit"`
	Credit    int64     `db:"credit"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

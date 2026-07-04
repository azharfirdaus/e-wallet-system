package model

import "time"

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyIDR Currency = "IDR"
	CurrencySGD Currency = "SGD"
	CurrencyJPY Currency = "JPY"
	CurrencyAUD Currency = "AUD"
	CurrencyCNY Currency = "CNY"
)

type WalletCurrency struct {
	ID                      int64     `db:"id"`
	WalletID                int64     `db:"wallet_id"`
	Currency                Currency  `db:"currency"`
	ClosingBalance          int64     `db:"closing_balance"`
	ClosingBalanceUpdatedAt time.Time `db:"closing_balance_updated_at"`
	CreatedAt               time.Time `db:"created_at"`
	UpdatedAt               time.Time `db:"updated_at"`
}

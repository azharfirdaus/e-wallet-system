package model

import "time"

type WalletStatus string

const (
	WalletStatusActivate  WalletStatus = "ACTIVATE"
	WalletStatusSuspended WalletStatus = "SUSPENDED"
)

type Wallet struct {
	ID                      int64        `db:"id"`
	UserID                  int64        `db:"user_id"`
	Status                  WalletStatus `db:"status"`
	ClosingBalance          int64        `db:"closing_balance"`
	ClosingBalanceUpdatedAt time.Time    `db:"closing_balance_updated_at"`
	CreatedAt               time.Time    `db:"created_at"`
	UpdatedAt               time.Time    `db:"updated_at"`
}

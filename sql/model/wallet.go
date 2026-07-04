package model

import "time"

type WalletStatus string

const (
	WalletStatusActivate  WalletStatus = "ACTIVATE"
	WalletStatusSuspended WalletStatus = "SUSPENDED"
)

type Wallet struct {
	ID        int64        `db:"id"`
	UserID    int64        `db:"user_id"`
	Status    WalletStatus `db:"status"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}

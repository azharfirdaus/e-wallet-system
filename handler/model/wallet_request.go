package model

type CreateWalletRequest struct {
	UserID       int64  `json:"user_id"`
	CurrencyCode string `json:"currency_code"`
}

type TopUpWalletRequest struct {
	CurrencyCode string `json:"currency_code"`
	Amount       string `json:"amount"`
}

type PayWithWalletRequest struct {
	CurrencyCode string `json:"currency_code"`
	Amount       string `json:"amount"`
}

type TransferWalletRequest struct {
	CurrencyCode string `json:"currency_code"`
	ToWalletID   int64  `json:"to_wallet_id"`
	Amount       string `json:"amount"`
}

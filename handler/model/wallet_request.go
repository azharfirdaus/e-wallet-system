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
	WalletIDDestination     int64  `json:"wallet_id_destination"`
	CurrencyCodeDestination string `json:"currency_code_destination"`
	Amount                  string `json:"amount"`
}

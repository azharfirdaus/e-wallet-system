package model

type CreateWalletResponse struct {
	WalletID     int64  `json:"wallet_id"`
	CurrencyCode string `json:"currency_code"`
}

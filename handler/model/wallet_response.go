package model

type CreateWalletResponse struct {
	WalletID     int64  `json:"wallet_id"`
	CurrencyCode string `json:"currency_code"`
}

type GetWalletResponse struct {
	WalletCurrencyID int64  `json:"wallet_currency_id"`
	CurrencyCode     string `json:"currency_code"`
	Balance          string `json:"balance"`
}

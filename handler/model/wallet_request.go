package model

type CreateWalletRequest struct {
	UserID       int64  `json:"user_id"`
	CurrencyCode string `json:"currency_code"`
}

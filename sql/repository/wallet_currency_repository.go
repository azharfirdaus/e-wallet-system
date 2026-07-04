package repository

import (
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/jmoiron/sqlx"
)

type WalletCurrencyRepository interface {
	Insert(trx *sqlx.Tx, m *model.WalletCurrency) (*int64, error)
}

type WalletCurrencyRepositoryImpl struct{}

func NewWalletCurrencyRepository() *WalletCurrencyRepositoryImpl {
	return &WalletCurrencyRepositoryImpl{}
}

func (w *WalletCurrencyRepositoryImpl) Insert(trx *sqlx.Tx, m *model.WalletCurrency) (*int64, error) {
	query := `
		INSERT INTO public.wallets_currency (wallet_id, currency)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	if err := trx.QueryRowx(query, m.WalletID, m.Currency).Scan(&id); err != nil {
		return nil, err
	}

	return &id, nil
}

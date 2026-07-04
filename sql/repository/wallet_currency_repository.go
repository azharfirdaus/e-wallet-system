package repository

import (
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/jmoiron/sqlx"
)

type WalletCurrencyRepository interface {
	Insert(m *model.WalletCurrency) (*int64, error)
}

type WalletCurrencyRepositoryImpl struct {
	db *sqlx.DB
}

func NewWalletCurrencyRepository(db *sqlx.DB) *WalletCurrencyRepositoryImpl {
	return &WalletCurrencyRepositoryImpl{db: db}
}

func (w *WalletCurrencyRepositoryImpl) Insert(m *model.WalletCurrency) (*int64, error) {
	trx, err := w.db.Beginx()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO public.wallets_currency (wallet_id, currency)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	if err := trx.QueryRowx(query, m.WalletID, m.Currency).Scan(&id); err != nil {
		_ = trx.Rollback()
		return nil, err
	}

	if err := trx.Commit(); err != nil {
		return nil, err
	}

	return &id, nil
}

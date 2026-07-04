package repository

import (
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/jmoiron/sqlx"
)

type WalletCurrencyRepository interface {
	Insert(trx *sqlx.Tx, m *model.WalletCurrency) (*int64, error)
	FindByWalletID(trx *sqlx.Tx, walletID int64) ([]model.WalletCurrency, error)
	FindByWalletIDAndCurrency(trx *sqlx.Tx, walletID int64, currency model.Currency) (*model.WalletCurrency, error)
}

type WalletCurrencyRepositoryImpl struct{}

func NewWalletCurrencyRepository() *WalletCurrencyRepositoryImpl {
	return &WalletCurrencyRepositoryImpl{}
}

func (w *WalletCurrencyRepositoryImpl) Insert(trx *sqlx.Tx, m *model.WalletCurrency) (*int64, error) {
	query := `
		INSERT INTO public.wallets_currency (wallet_id, currency)
		VALUES ($1, $2)
		RETURNING id, closing_balance, closing_balance_updated_at
	`

	var id int64
	if err := trx.QueryRowx(query, m.WalletID, m.Currency).Scan(
		&id,
		&m.ClosingBalance,
		&m.ClosingBalanceUpdatedAt,
	); err != nil {
		return nil, err
	}

	return &id, nil
}

func (w *WalletCurrencyRepositoryImpl) FindByWalletID(
	trx *sqlx.Tx,
	walletID int64,
) ([]model.WalletCurrency, error) {
	query := `
		SELECT
			id,
			wallet_id,
			currency,
			closing_balance,
			closing_balance_updated_at,
			created_at,
			updated_at
		FROM public.wallets_currency
		WHERE wallet_id = $1
		ORDER BY id
	`

	walletCurrencies := []model.WalletCurrency{}
	if err := trx.Select(&walletCurrencies, query, walletID); err != nil {
		return nil, err
	}

	return walletCurrencies, nil
}

func (w *WalletCurrencyRepositoryImpl) FindByWalletIDAndCurrency(
	trx *sqlx.Tx,
	walletID int64,
	currency model.Currency,
) (*model.WalletCurrency, error) {
	query := `
		SELECT
			id,
			wallet_id,
			currency,
			closing_balance,
			closing_balance_updated_at,
			created_at,
			updated_at
		FROM public.wallets_currency
		WHERE wallet_id = $1
			AND currency = $2
	`

	walletCurrency := model.WalletCurrency{}
	if err := trx.QueryRowx(query, walletID, currency).StructScan(&walletCurrency); err != nil {
		return nil, err
	}

	return &walletCurrency, nil
}

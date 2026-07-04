package repository

import (
	"database/sql"

	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/jmoiron/sqlx"
)

type WalletRepository interface {
	Insert(trx *sqlx.Tx, m *model.Wallet) (*int64, error)
	FindByID(trx *sqlx.Tx, id int64) (*model.Wallet, error)
	Suspend(trx *sqlx.Tx, id int64) error
}

type WalletRepositoryImpl struct{}

func NewWalletRepository() *WalletRepositoryImpl {
	return &WalletRepositoryImpl{}
}

func (w *WalletRepositoryImpl) Insert(trx *sqlx.Tx, m *model.Wallet) (*int64, error) {
	query := `
		INSERT INTO public.wallets (user_id, status, closing_balance)
		VALUES ($1, $2, $3)
		RETURNING id, closing_balance, closing_balance_updated_at
	`

	var id int64
	if err := trx.QueryRowx(query, m.UserID, m.Status, m.ClosingBalance).Scan(
		&id,
		&m.ClosingBalance,
		&m.ClosingBalanceUpdatedAt,
	); err != nil {
		return nil, err
	}

	return &id, nil
}

func (w *WalletRepositoryImpl) FindByID(trx *sqlx.Tx, id int64) (*model.Wallet, error) {
	query := `
		SELECT
			id,
			user_id,
			status,
			closing_balance,
			closing_balance_updated_at,
			created_at,
			updated_at
		FROM public.wallets
		WHERE id = $1
	`

	wallet := model.Wallet{}
	if err := trx.QueryRowx(query, id).StructScan(&wallet); err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (w *WalletRepositoryImpl) Suspend(trx *sqlx.Tx, id int64) error {
	result, err := trx.Exec(
		`UPDATE public.wallets
		 SET status = $1
		 WHERE id = $2`,
		model.WalletStatusSuspended,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

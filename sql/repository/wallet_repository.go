package repository

import (
	"database/sql"

	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/jmoiron/sqlx"
)

type WalletRepository interface {
	Insert(trx *sqlx.Tx, m *model.Wallet) (*int64, error)
	FindByID(trx *sqlx.Tx, id int64) (*model.Wallet, error)
	FindByUserID(trx *sqlx.Tx, userID int64) (*model.Wallet, error)
	Suspend(trx *sqlx.Tx, id int64) error
}

type WalletRepositoryImpl struct{}

func NewWalletRepository() *WalletRepositoryImpl {
	return &WalletRepositoryImpl{}
}

func (w *WalletRepositoryImpl) Insert(trx *sqlx.Tx, m *model.Wallet) (*int64, error) {
	query := `
		INSERT INTO public.wallets (user_id, status)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	if err := trx.QueryRowx(query, m.UserID, m.Status).Scan(&id); err != nil {
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

func (w *WalletRepositoryImpl) FindByUserID(trx *sqlx.Tx, userID int64) (*model.Wallet, error) {
	query := `
		SELECT
			id,
			user_id,
			status,
			created_at,
			updated_at
		FROM public.wallets
		WHERE user_id = $1
		ORDER BY id
		LIMIT 1
	`

	wallet := model.Wallet{}
	if err := trx.QueryRowx(query, userID).StructScan(&wallet); err != nil {
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

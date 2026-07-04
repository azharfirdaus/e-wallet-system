package repository

import (
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/jmoiron/sqlx"
)

type WalletRepository interface {
	Insert(m *model.Wallet) (*int64, error)
}

type WalletRepositoryImpl struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) *WalletRepositoryImpl {
	return &WalletRepositoryImpl{db: db}
}

func (w *WalletRepositoryImpl) Insert(m *model.Wallet) (*int64, error) {
	trx, err := w.db.Beginx()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO public.wallets (user_id, status)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	if err := trx.QueryRowx(query, m.UserID, m.Status).Scan(&id); err != nil {
		_ = trx.Rollback()
		return nil, err
	}

	if err := trx.Commit(); err != nil {
		return nil, err
	}

	return &id, nil
}

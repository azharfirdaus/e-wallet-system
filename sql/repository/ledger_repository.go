package repository

import (
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/jmoiron/sqlx"
)

type LedgerRepository interface {
	Insert(m *model.Ledger) (*int64, error)
}

type LedgerRepositoryImpl struct {
	db *sqlx.DB
}

func NewLedgerRepository(db *sqlx.DB) *LedgerRepositoryImpl {
	return &LedgerRepositoryImpl{db: db}
}

func (l *LedgerRepositoryImpl) Insert(m *model.Ledger) (*int64, error) {
	trx, err := l.db.Beginx()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO public.ledger (debit, credit)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	if err := trx.QueryRowx(query, m.Debit, m.Credit).Scan(&id); err != nil {
		_ = trx.Rollback()
		return nil, err
	}

	if err := trx.Commit(); err != nil {
		return nil, err
	}

	return &id, nil
}

package repository

import (
	"time"

	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/jmoiron/sqlx"
)

type LedgerRepository interface {
	Insert(trx *sqlx.Tx, m *model.Ledger) (*int64, error)
	SumByWalletCurrencyID(trx *sqlx.Tx, afterCreatedAt time.Time) ([]model.LedgerSum, error)
}

type LedgerRepositoryImpl struct{}

func NewLedgerRepository() *LedgerRepositoryImpl {
	return &LedgerRepositoryImpl{}
}

func (l *LedgerRepositoryImpl) Insert(trx *sqlx.Tx, m *model.Ledger) (*int64, error) {
	query := `
		INSERT INTO public.ledger (wallet_currency_id, debit, credit, reference)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	var id int64
	if err := trx.QueryRowx(query, m.WalletCurrencyID, m.Debit, m.Credit, m.Reference).Scan(
		&id,
		&m.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &id, nil
}

func (l *LedgerRepositoryImpl) SumByWalletCurrencyID(
	trx *sqlx.Tx,
	afterCreatedAt time.Time,
) ([]model.LedgerSum, error) {
	query := `
		SELECT
			wallet_currency_id,
			COALESCE(SUM(debit), 0) AS debit,
			COALESCE(SUM(credit), 0) AS credit
		FROM public.ledger
		WHERE created_at > $1
		GROUP BY wallet_currency_id
	`

	sums := []model.LedgerSum{}
	if err := trx.Select(&sums, query, afterCreatedAt); err != nil {
		return nil, err
	}

	return sums, nil
}

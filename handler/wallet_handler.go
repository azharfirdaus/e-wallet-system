package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/azharfirdaus/e-wallet-system/sql/repository"
	"github.com/jmoiron/sqlx"
)

type WalletHandler struct {
	db                       *sqlx.DB
	walletRepository         repository.WalletRepository
	walletCurrencyRepository repository.WalletCurrencyRepository
	ledgerRepository         repository.LedgerRepository
}

const walletWriteAdvisoryLockID int64 = 2026070401

var amountLimitByCurrency = map[model.Currency]int64{
	model.CurrencyIDR: 1_000_000_000, // 10,000,000.00
	model.CurrencyUSD: 60_000,        // 600.00
	model.CurrencySGD: 78_000,        // 780.00
	model.CurrencyJPY: 9_000_000,     // 90,000.00
	model.CurrencyAUD: 93_000,        // 930.00
	model.CurrencyCNY: 440_000,       // 4,400.00
}

func NewWalletHandler(db *sqlx.DB) *WalletHandler {
	return &WalletHandler{
		db:                       db,
		walletRepository:         repository.NewWalletRepository(),
		walletCurrencyRepository: repository.NewWalletCurrencyRepository(),
		ledgerRepository:         repository.NewLedgerRepository(),
	}
}

func validateAmountLimit(amount int64, currency model.Currency) error {
	limit, ok := amountLimitByCurrency[currency]
	if !ok {
		return errors.New("unsupported currency_code")
	}
	if amount > limit {
		return fmt.Errorf("amount exceeds %s limit", currency)
	}

	return nil
}

func acquireWalletWriteLock(trx *sqlx.Tx) error {
	_, err := trx.Exec("SELECT pg_advisory_xact_lock($1)", walletWriteAdvisoryLockID)
	return err
}

func parseCurrency(currencyCode string) (model.Currency, error) {
	currency := model.Currency(strings.ToUpper(currencyCode))
	switch currency {
	case model.CurrencyUSD,
		model.CurrencyIDR,
		model.CurrencySGD,
		model.CurrencyJPY,
		model.CurrencyAUD,
		model.CurrencyCNY:
		return currency, nil
	default:
		return "", errors.New("unsupported currency_code")
	}
}

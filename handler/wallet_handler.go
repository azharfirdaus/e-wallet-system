package handler

import (
	"errors"
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

func NewWalletHandler(db *sqlx.DB) *WalletHandler {
	return &WalletHandler{
		db:                       db,
		walletRepository:         repository.NewWalletRepository(),
		walletCurrencyRepository: repository.NewWalletCurrencyRepository(),
		ledgerRepository:         repository.NewLedgerRepository(),
	}
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

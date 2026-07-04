package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	handlermodel "github.com/azharfirdaus/e-wallet-system/handler/model"
	"github.com/azharfirdaus/e-wallet-system/helper"
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/azharfirdaus/e-wallet-system/sql/repository"
	"github.com/gorilla/mux"
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

func (h *WalletHandler) CreateWalletHandler(w http.ResponseWriter, r *http.Request) {
	var request handlermodel.CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.UserID <= 0 {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	currency, err := parseCurrency(request.CurrencyCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	trx, err := h.db.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	walletID, err := h.walletRepository.Insert(trx, &model.Wallet{
		UserID: request.UserID,
		Status: model.WalletStatusActivate,
	})
	if err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := h.walletCurrencyRepository.Insert(trx, &model.WalletCurrency{
		WalletID: *walletID,
		Currency: currency,
	}); err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := trx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(handlermodel.CreateWalletResponse{
		WalletID:     *walletID,
		CurrencyCode: string(currency),
	})
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

func (h *WalletHandler) TopUpWalletHandler(w http.ResponseWriter, r *http.Request) {
	walletID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil || walletID <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var request handlermodel.TopUpWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	amount, err := helper.ParseAmountMinorUnit(request.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	currency, err := parseCurrency(request.CurrencyCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	trx, err := h.db.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	walletCurrency, err := h.walletCurrencyRepository.FindByWalletIDAndCurrency(trx, walletID, currency)
	if err != nil {
		_ = trx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "wallet currency not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := h.ledgerRepository.Insert(trx, &model.Ledger{
		WalletCurrencyID: walletCurrency.ID,
		Debit:            0,
		Credit:           amount,
		Reference:        model.LedgerReferenceTopup,
	}); err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := trx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WalletHandler) PayWithWalletHandler(w http.ResponseWriter, r *http.Request) {
	walletID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil || walletID <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var request handlermodel.PayWithWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	amount, err := helper.ParseAmountMinorUnit(request.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	currency, err := parseCurrency(request.CurrencyCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	trx, err := h.db.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	walletCurrency, err := h.walletCurrencyRepository.FindByWalletIDAndCurrency(trx, walletID, currency)
	if err != nil {
		_ = trx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "wallet currency not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := h.ledgerRepository.Insert(trx, &model.Ledger{
		WalletCurrencyID: walletCurrency.ID,
		Debit:            amount,
		Credit:           0,
		Reference:        model.LedgerReferencePayment,
	}); err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := trx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func TransferWalletHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *WalletHandler) SuspendWalletHandler(w http.ResponseWriter, r *http.Request) {
	walletID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil || walletID <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	trx, err := h.db.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.walletRepository.Suspend(trx, walletID); err != nil {
		_ = trx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "wallet not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := trx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WalletHandler) GetWalletHandler(w http.ResponseWriter, r *http.Request) {
	walletID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil || walletID <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	trx, err := h.db.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wallet, err := h.walletRepository.FindByID(trx, walletID)
	if err != nil {
		_ = trx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "wallet not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	walletCurrencies, err := h.walletCurrencyRepository.FindByWalletID(trx, walletID)
	if err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(walletCurrencies) == 0 {
		_ = trx.Rollback()
		http.Error(w, "wallet currency not found", http.StatusNotFound)
		return
	}

	ledgerSums, err := h.ledgerRepository.SumByWalletCurrencyID(trx, wallet.ClosingBalanceUpdatedAt)
	if err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := trx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ledgerSumByWalletCurrencyID := make(map[int64]model.LedgerSum, len(ledgerSums))
	for _, ledgerSum := range ledgerSums {
		ledgerSumByWalletCurrencyID[ledgerSum.WalletCurrencyID] = ledgerSum
	}

	response := make([]handlermodel.GetWalletResponse, 0, len(walletCurrencies))
	for _, walletCurrency := range walletCurrencies {
		ledgerSum := ledgerSumByWalletCurrencyID[walletCurrency.ID]
		response = append(response, handlermodel.GetWalletResponse{
			WalletCurrencyID: walletCurrency.ID,
			CurrencyCode:     string(walletCurrency.Currency),
			Balance:          ledgerSum.Credit - ledgerSum.Debit,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

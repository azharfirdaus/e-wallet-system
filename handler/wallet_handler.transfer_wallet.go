package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	handlermodel "github.com/azharfirdaus/e-wallet-system/handler/model"
	"github.com/azharfirdaus/e-wallet-system/helper"
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/gorilla/mux"
)

func (h *WalletHandler) TransferWalletHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || walletID <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var request handlermodel.TransferWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.WalletIDDestination <= 0 {
		http.Error(w, "wallet_id_destination is required", http.StatusBadRequest)
		return
	}
	if request.WalletIDDestination == walletID {
		http.Error(w, "wallet_id_destination must be different from wallet id", http.StatusBadRequest)
		return
	}

	trx, err := h.db.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := acquireWalletWriteLock(trx); err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fromWallet, err := h.walletRepository.FindByID(trx, walletID)
	if err != nil {
		_ = trx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "wallet not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if fromWallet.Status != model.WalletStatusActivate {
		_ = trx.Rollback()
		http.Error(w, "wallet is not activated", http.StatusUnprocessableEntity)
		return
	}

	toWallet, err := h.walletRepository.FindByID(trx, request.WalletIDDestination)
	if err != nil {
		_ = trx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "destination wallet not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if toWallet.Status != model.WalletStatusActivate {
		_ = trx.Rollback()
		http.Error(w, "destination wallet is not activated", http.StatusUnprocessableEntity)
		return
	}

	amount, err := helper.ParseAmountMinorUnit(request.Amount)
	if err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	sourceCurrency, err := parseCurrency(vars["country_code"])
	if err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	destinationCurrency, err := parseCurrency(request.CurrencyCodeDestination)
	if err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if sourceCurrency != destinationCurrency {
		_ = trx.Rollback()
		http.Error(w, "currency code source and destination are different", http.StatusUnprocessableEntity)
		return
	}
	if err := validateAmountLimit(amount, sourceCurrency); err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	fromWalletCurrency, err := h.walletCurrencyRepository.FindByWalletIDAndCurrency(trx, walletID, sourceCurrency)
	if err != nil {
		_ = trx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "wallet currency not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	toWalletCurrency, err := h.walletCurrencyRepository.FindByWalletIDAndCurrency(trx, request.WalletIDDestination, destinationCurrency)
	if err != nil {
		_ = trx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "destination wallet currency not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ledgerSum, err := h.ledgerRepository.SumOneByWalletCurrencyID(
		trx,
		fromWalletCurrency.ID,
		fromWalletCurrency.ClosingBalanceUpdatedAt,
	)
	if err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if fromWalletCurrency.ClosingBalance+ledgerSum.Credit-ledgerSum.Debit < amount {
		_ = trx.Rollback()
		http.Error(w, "insufficient amount", http.StatusUnprocessableEntity)
		return
	}

	if _, err := h.ledgerRepository.Insert(trx, &model.Ledger{
		WalletCurrencyID: fromWalletCurrency.ID,
		Debit:            amount,
		Credit:           0,
		Reference:        model.LedgerReferenceTransfer,
	}); err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := h.ledgerRepository.Insert(trx, &model.Ledger{
		WalletCurrencyID: toWalletCurrency.ID,
		Debit:            0,
		Credit:           amount,
		Reference:        model.LedgerReferenceReceive,
	}); err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := trx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

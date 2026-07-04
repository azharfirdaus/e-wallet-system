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
	walletID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil || walletID <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var request handlermodel.TransferWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.ToWalletID <= 0 {
		http.Error(w, "to_wallet_id is required", http.StatusBadRequest)
		return
	}
	if request.ToWalletID == walletID {
		http.Error(w, "to_wallet_id must be different from wallet id", http.StatusBadRequest)
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
	if err := validateAmountLimit(amount, currency); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	trx, err := h.db.Beginx()
	if err != nil {
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

	toWallet, err := h.walletRepository.FindByID(trx, request.ToWalletID)
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

	fromWalletCurrency, err := h.walletCurrencyRepository.FindByWalletIDAndCurrency(trx, walletID, currency)
	if err != nil {
		_ = trx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "wallet currency not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	toWalletCurrency, err := h.walletCurrencyRepository.FindByWalletIDAndCurrency(trx, request.ToWalletID, currency)
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
		fromWallet.ClosingBalanceUpdatedAt,
	)
	if err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ledgerSum.Credit-ledgerSum.Debit < amount {
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

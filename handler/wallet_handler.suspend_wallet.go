package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/gorilla/mux"
)

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
	if err := acquireWalletWriteLock(trx); err != nil {
		_ = trx.Rollback()
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
	if wallet.Status == model.WalletStatusSuspended {
		_ = trx.Rollback()
		http.Error(w, "wallet is already suspended", http.StatusUnprocessableEntity)
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

	w.WriteHeader(http.StatusOK)
}

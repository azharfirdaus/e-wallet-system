package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

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

package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	handlermodel "github.com/azharfirdaus/e-wallet-system/handler/model"
	"github.com/gorilla/mux"
)

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

	response := make([]handlermodel.GetWalletResponse, 0, len(walletCurrencies))
	for _, walletCurrency := range walletCurrencies {
		ledgerSum, err := h.ledgerRepository.SumOneByWalletCurrencyID(
			trx,
			walletCurrency.ID,
			walletCurrency.ClosingBalanceUpdatedAt,
		)
		if err != nil {
			_ = trx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		balance := walletCurrency.ClosingBalance + ledgerSum.Credit - ledgerSum.Debit
		response = append(response, handlermodel.GetWalletResponse{
			WalletCurrencyID: walletCurrency.ID,
			CurrencyCode:     string(walletCurrency.Currency),
			Status:           string(wallet.Status),
			Balance:          formatBalance(balance),
		})
	}

	if err := trx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func formatBalance(balance int64) string {
	if balance < 0 {
		return fmt.Sprintf("-%d.%02d", (-balance)/100, (-balance)%100)
	}

	return fmt.Sprintf("%d.%02d", balance/100, balance%100)
}

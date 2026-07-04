package handler

import (
	"encoding/json"
	"net/http"

	handlermodel "github.com/azharfirdaus/e-wallet-system/handler/model"
	"github.com/azharfirdaus/e-wallet-system/sql/model"
)

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

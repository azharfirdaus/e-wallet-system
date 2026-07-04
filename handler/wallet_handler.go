package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	handlermodel "github.com/azharfirdaus/e-wallet-system/handler/model"
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/azharfirdaus/e-wallet-system/sql/repository"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type WalletHandler struct {
	db *sqlx.DB
}

func NewWalletHandler(db *sqlx.DB) *WalletHandler {
	return &WalletHandler{db: db}
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

	walletRepository := repository.NewWalletRepository()
	walletID, err := walletRepository.Insert(trx, &model.Wallet{
		UserID: request.UserID,
		Status: model.WalletStatusActivate,
	})
	if err != nil {
		_ = trx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	walletCurrencyRepository := repository.NewWalletCurrencyRepository()
	if _, err := walletCurrencyRepository.Insert(trx, &model.WalletCurrency{
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

func TopUpWalletHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func PayWithWalletHandler(w http.ResponseWriter, _ *http.Request) {
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

	walletRepository := repository.NewWalletRepository()
	if err := walletRepository.Suspend(trx, walletID); err != nil {
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

func GetWalletHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

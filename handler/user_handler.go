package handler

import (
	"encoding/json"
	"net/http"

	handlermodel "github.com/azharfirdaus/e-wallet-system/handler/model"
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/azharfirdaus/e-wallet-system/sql/repository"
	"github.com/jmoiron/sqlx"
)

type UserHandler struct {
	db             *sqlx.DB
	userRepository repository.UserRepository
}

func NewUserHandler(db *sqlx.DB) *UserHandler {
	return &UserHandler{
		db:             db,
		userRepository: repository.NewUserRepository(),
	}
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, _ *http.Request) {
	trx, err := h.db.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := h.userRepository.Insert(trx, &model.User{})
	if err != nil {
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
	_ = json.NewEncoder(w).Encode(handlermodel.CreateUserResponse{
		ID: *id,
	})
}

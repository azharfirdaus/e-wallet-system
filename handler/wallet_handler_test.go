package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	handlermodel "github.com/azharfirdaus/e-wallet-system/handler/model"
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type mockWalletRepository struct {
	insertFunc       func(trx *sqlx.Tx, wallet *model.Wallet) (*int64, error)
	findByIDFunc     func(trx *sqlx.Tx, id int64) (*model.Wallet, error)
	findByUserIDFunc func(trx *sqlx.Tx, userID int64) (*model.Wallet, error)
	suspendFunc      func(trx *sqlx.Tx, id int64) error
}

func (m *mockWalletRepository) Insert(trx *sqlx.Tx, wallet *model.Wallet) (*int64, error) {
	if m.insertFunc == nil {
		panic("mockWalletRepository.Insert is not implemented")
	}
	return m.insertFunc(trx, wallet)
}

func (m *mockWalletRepository) FindByID(trx *sqlx.Tx, id int64) (*model.Wallet, error) {
	if m.findByIDFunc == nil {
		panic("mockWalletRepository.FindByID is not implemented")
	}
	return m.findByIDFunc(trx, id)
}

func (m *mockWalletRepository) FindByUserID(trx *sqlx.Tx, userID int64) (*model.Wallet, error) {
	if m.findByUserIDFunc == nil {
		panic("mockWalletRepository.FindByUserID is not implemented")
	}
	return m.findByUserIDFunc(trx, userID)
}

func (m *mockWalletRepository) Suspend(trx *sqlx.Tx, id int64) error {
	if m.suspendFunc == nil {
		panic("mockWalletRepository.Suspend is not implemented")
	}
	return m.suspendFunc(trx, id)
}

type mockWalletCurrencyRepository struct {
	insertFunc                    func(trx *sqlx.Tx, walletCurrency *model.WalletCurrency) (*int64, error)
	findByWalletIDFunc            func(trx *sqlx.Tx, walletID int64) ([]model.WalletCurrency, error)
	findByWalletIDAndCurrencyFunc func(trx *sqlx.Tx, walletID int64, currency model.Currency) (*model.WalletCurrency, error)
	updateClosingBalanceFunc      func(trx *sqlx.Tx) error
}

func (m *mockWalletCurrencyRepository) Insert(
	trx *sqlx.Tx,
	walletCurrency *model.WalletCurrency,
) (*int64, error) {
	if m.insertFunc == nil {
		panic("mockWalletCurrencyRepository.Insert is not implemented")
	}
	return m.insertFunc(trx, walletCurrency)
}

func (m *mockWalletCurrencyRepository) FindByWalletID(
	trx *sqlx.Tx,
	walletID int64,
) ([]model.WalletCurrency, error) {
	if m.findByWalletIDFunc == nil {
		panic("mockWalletCurrencyRepository.FindByWalletID is not implemented")
	}
	return m.findByWalletIDFunc(trx, walletID)
}

func (m *mockWalletCurrencyRepository) FindByWalletIDAndCurrency(
	trx *sqlx.Tx,
	walletID int64,
	currency model.Currency,
) (*model.WalletCurrency, error) {
	if m.findByWalletIDAndCurrencyFunc == nil {
		panic("mockWalletCurrencyRepository.FindByWalletIDAndCurrency is not implemented")
	}
	return m.findByWalletIDAndCurrencyFunc(trx, walletID, currency)
}

func (m *mockWalletCurrencyRepository) UpdateClosingBalanceByLedger(trx *sqlx.Tx) error {
	if m.updateClosingBalanceFunc == nil {
		panic("mockWalletCurrencyRepository.UpdateClosingBalanceByLedger is not implemented")
	}
	return m.updateClosingBalanceFunc(trx)
}

type mockLedgerRepository struct {
	insertFunc                 func(trx *sqlx.Tx, ledger *model.Ledger) (*int64, error)
	sumByWalletCurrencyIDFunc  func(trx *sqlx.Tx, afterCreatedAt time.Time) ([]model.LedgerSum, error)
	sumOneByWalletCurrencyFunc func(trx *sqlx.Tx, walletCurrencyID int64, afterCreatedAt time.Time) (*model.LedgerSum, error)
}

func (m *mockLedgerRepository) Insert(trx *sqlx.Tx, ledger *model.Ledger) (*int64, error) {
	if m.insertFunc == nil {
		panic("mockLedgerRepository.Insert is not implemented")
	}
	return m.insertFunc(trx, ledger)
}

func (m *mockLedgerRepository) SumByWalletCurrencyID(
	trx *sqlx.Tx,
	afterCreatedAt time.Time,
) ([]model.LedgerSum, error) {
	if m.sumByWalletCurrencyIDFunc == nil {
		panic("mockLedgerRepository.SumByWalletCurrencyID is not implemented")
	}
	return m.sumByWalletCurrencyIDFunc(trx, afterCreatedAt)
}

func (m *mockLedgerRepository) SumOneByWalletCurrencyID(
	trx *sqlx.Tx,
	walletCurrencyID int64,
	afterCreatedAt time.Time,
) (*model.LedgerSum, error) {
	if m.sumOneByWalletCurrencyFunc == nil {
		panic("mockLedgerRepository.SumOneByWalletCurrencyID is not implemented")
	}
	return m.sumOneByWalletCurrencyFunc(trx, walletCurrencyID, afterCreatedAt)
}

func expectWalletWriteLock(sqlMock sqlmock.Sqlmock) {
	sqlMock.ExpectExec("SELECT pg_advisory_xact_lock\\(\\$1\\)").
		WithArgs(walletWriteAdvisoryLockID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func requestWithVars(method string, target string, body string, vars map[string]string) *http.Request {
	request := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	return mux.SetURLVars(request, vars)
}

func TestCreateWalletHandler(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		db, sqlMock, cleanup := newHandlerTestDB(t)
		defer cleanup()

		walletID := int64(101)
		walletCurrencyID := int64(202)
		sqlMock.ExpectBegin()
		expectWalletWriteLock(sqlMock)
		sqlMock.ExpectCommit()

		handler := &WalletHandler{
			db: db,
			walletRepository: &mockWalletRepository{
				findByUserIDFunc: func(trx *sqlx.Tx, userID int64) (*model.Wallet, error) {
					if trx == nil {
						t.Fatal("expected transaction")
					}
					if userID != 7 {
						t.Fatalf("expected user id 7, got %d", userID)
					}
					return nil, sql.ErrNoRows
				},
				insertFunc: func(trx *sqlx.Tx, wallet *model.Wallet) (*int64, error) {
					if wallet.UserID != 7 {
						t.Fatalf("expected user id 7, got %d", wallet.UserID)
					}
					if wallet.Status != model.WalletStatusActivate {
						t.Fatalf("expected status %s, got %s", model.WalletStatusActivate, wallet.Status)
					}
					return &walletID, nil
				},
			},
			walletCurrencyRepository: &mockWalletCurrencyRepository{
				insertFunc: func(trx *sqlx.Tx, walletCurrency *model.WalletCurrency) (*int64, error) {
					if walletCurrency.WalletID != walletID {
						t.Fatalf("expected wallet id %d, got %d", walletID, walletCurrency.WalletID)
					}
					if walletCurrency.Currency != model.CurrencyUSD {
						t.Fatalf("expected currency %s, got %s", model.CurrencyUSD, walletCurrency.Currency)
					}
					return &walletCurrencyID, nil
				},
			},
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(
			http.MethodPost,
			"/wallets",
			bytes.NewBufferString(`{"user_id":7,"currency_code":"usd"}`),
		)
		handler.CreateWalletHandler(recorder, request)

		if recorder.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
		}
		if recorder.Header().Get("Content-Type") != "application/json" {
			t.Fatalf("expected content type application/json, got %s", recorder.Header().Get("Content-Type"))
		}

		var response handlermodel.CreateWalletResponse
		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if response.WalletID != walletID {
			t.Fatalf("expected wallet id %d, got %d", walletID, response.WalletID)
		}
		if response.CurrencyCode != string(model.CurrencyUSD) {
			t.Fatalf("expected currency %s, got %s", model.CurrencyUSD, response.CurrencyCode)
		}
	})

	t.Run("repository error rolls back", func(t *testing.T) {
		db, sqlMock, cleanup := newHandlerTestDB(t)
		defer cleanup()

		sqlMock.ExpectBegin()
		expectWalletWriteLock(sqlMock)
		sqlMock.ExpectRollback()

		handler := &WalletHandler{
			db: db,
			walletRepository: &mockWalletRepository{
				findByUserIDFunc: func(trx *sqlx.Tx, userID int64) (*model.Wallet, error) {
					return nil, sql.ErrConnDone
				},
			},
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(
			http.MethodPost,
			"/wallets",
			bytes.NewBufferString(`{"user_id":7,"currency_code":"USD"}`),
		)
		handler.CreateWalletHandler(recorder, request)

		if recorder.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
		}
	})
}

func TestGetWalletHandler(t *testing.T) {
	db, sqlMock, cleanup := newHandlerTestDB(t)
	defer cleanup()

	updatedAt := time.Date(2026, 7, 5, 9, 0, 0, 0, time.UTC)
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	handler := &WalletHandler{
		db: db,
		walletRepository: &mockWalletRepository{
			findByIDFunc: func(trx *sqlx.Tx, id int64) (*model.Wallet, error) {
				if id != 10 {
					t.Fatalf("expected wallet id 10, got %d", id)
				}
				return &model.Wallet{ID: 10, UserID: 7, Status: model.WalletStatusActivate}, nil
			},
		},
		walletCurrencyRepository: &mockWalletCurrencyRepository{
			findByWalletIDFunc: func(trx *sqlx.Tx, walletID int64) ([]model.WalletCurrency, error) {
				if walletID != 10 {
					t.Fatalf("expected wallet id 10, got %d", walletID)
				}
				return []model.WalletCurrency{
					{
						ID:                      21,
						WalletID:                10,
						Currency:                model.CurrencyIDR,
						ClosingBalance:          10_000,
						ClosingBalanceUpdatedAt: updatedAt,
					},
				}, nil
			},
		},
		ledgerRepository: &mockLedgerRepository{
			sumOneByWalletCurrencyFunc: func(trx *sqlx.Tx, walletCurrencyID int64, afterCreatedAt time.Time) (*model.LedgerSum, error) {
				if walletCurrencyID != 21 {
					t.Fatalf("expected wallet currency id 21, got %d", walletCurrencyID)
				}
				if !afterCreatedAt.Equal(updatedAt) {
					t.Fatalf("expected after created at %s, got %s", updatedAt, afterCreatedAt)
				}
				return &model.LedgerSum{WalletCurrencyID: 21, Debit: 50, Credit: 250}, nil
			},
		},
	}

	recorder := httptest.NewRecorder()
	request := requestWithVars(http.MethodGet, "/wallets/10", "", map[string]string{"id": "10"})
	handler.GetWalletHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected content type application/json, got %s", recorder.Header().Get("Content-Type"))
	}

	var response []handlermodel.GetWalletResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	expected := handlermodel.GetWalletResponse{
		WalletCurrencyID: 21,
		CurrencyCode:     "IDR",
		Status:           "ACTIVATE",
		Balance:          "102.00",
	}
	if len(response) != 1 {
		t.Fatalf("expected 1 response item, got %d", len(response))
	}
	if response[0] != expected {
		t.Fatalf("expected response %#v, got %#v", expected, response[0])
	}
}

func TestTopUpWalletHandler(t *testing.T) {
	db, sqlMock, cleanup := newHandlerTestDB(t)
	defer cleanup()

	ledgerID := int64(301)
	sqlMock.ExpectBegin()
	expectWalletWriteLock(sqlMock)
	sqlMock.ExpectCommit()

	handler := &WalletHandler{
		db: db,
		walletRepository: &mockWalletRepository{
			findByIDFunc: func(trx *sqlx.Tx, id int64) (*model.Wallet, error) {
				if id != 10 {
					t.Fatalf("expected wallet id 10, got %d", id)
				}
				return &model.Wallet{ID: 10, UserID: 7, Status: model.WalletStatusActivate}, nil
			},
		},
		walletCurrencyRepository: &mockWalletCurrencyRepository{
			findByWalletIDAndCurrencyFunc: func(
				trx *sqlx.Tx,
				walletID int64,
				currency model.Currency,
			) (*model.WalletCurrency, error) {
				if walletID != 10 {
					t.Fatalf("expected wallet id 10, got %d", walletID)
				}
				if currency != model.CurrencyIDR {
					t.Fatalf("expected currency %s, got %s", model.CurrencyIDR, currency)
				}
				return &model.WalletCurrency{ID: 21, WalletID: 10, Currency: model.CurrencyIDR}, nil
			},
		},
		ledgerRepository: &mockLedgerRepository{
			insertFunc: func(trx *sqlx.Tx, ledger *model.Ledger) (*int64, error) {
				if ledger.WalletCurrencyID != 21 {
					t.Fatalf("expected wallet currency id 21, got %d", ledger.WalletCurrencyID)
				}
				if ledger.Debit != 0 || ledger.Credit != 12_345 {
					t.Fatalf("expected debit 0 and credit 12345, got debit %d credit %d", ledger.Debit, ledger.Credit)
				}
				if ledger.Reference != model.LedgerReferenceTopup {
					t.Fatalf("expected reference %s, got %s", model.LedgerReferenceTopup, ledger.Reference)
				}
				return &ledgerID, nil
			},
		},
	}

	recorder := httptest.NewRecorder()
	request := requestWithVars(
		http.MethodPost,
		"/wallets/10/topup",
		`{"currency_code":"IDR","amount":"123.45"}`,
		map[string]string{"id": "10"},
	)
	handler.TopUpWalletHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
}

func TestPayWithWalletHandler(t *testing.T) {
	db, sqlMock, cleanup := newHandlerTestDB(t)
	defer cleanup()

	updatedAt := time.Date(2026, 7, 5, 10, 0, 0, 0, time.UTC)
	ledgerID := int64(401)
	sqlMock.ExpectBegin()
	expectWalletWriteLock(sqlMock)
	sqlMock.ExpectCommit()

	handler := &WalletHandler{
		db: db,
		walletRepository: &mockWalletRepository{
			findByIDFunc: func(trx *sqlx.Tx, id int64) (*model.Wallet, error) {
				if id != 10 {
					t.Fatalf("expected wallet id 10, got %d", id)
				}
				return &model.Wallet{ID: 10, UserID: 7, Status: model.WalletStatusActivate}, nil
			},
		},
		walletCurrencyRepository: &mockWalletCurrencyRepository{
			findByWalletIDAndCurrencyFunc: func(
				trx *sqlx.Tx,
				walletID int64,
				currency model.Currency,
			) (*model.WalletCurrency, error) {
				if walletID != 10 {
					t.Fatalf("expected wallet id 10, got %d", walletID)
				}
				if currency != model.CurrencyUSD {
					t.Fatalf("expected currency %s, got %s", model.CurrencyUSD, currency)
				}
				return &model.WalletCurrency{
					ID:                      31,
					WalletID:                10,
					Currency:                model.CurrencyUSD,
					ClosingBalance:          10_000,
					ClosingBalanceUpdatedAt: updatedAt,
				}, nil
			},
		},
		ledgerRepository: &mockLedgerRepository{
			sumOneByWalletCurrencyFunc: func(trx *sqlx.Tx, walletCurrencyID int64, afterCreatedAt time.Time) (*model.LedgerSum, error) {
				if walletCurrencyID != 31 {
					t.Fatalf("expected wallet currency id 31, got %d", walletCurrencyID)
				}
				if !afterCreatedAt.Equal(updatedAt) {
					t.Fatalf("expected after created at %s, got %s", updatedAt, afterCreatedAt)
				}
				return &model.LedgerSum{WalletCurrencyID: 31, Debit: 1_000, Credit: 500}, nil
			},
			insertFunc: func(trx *sqlx.Tx, ledger *model.Ledger) (*int64, error) {
				if ledger.WalletCurrencyID != 31 {
					t.Fatalf("expected wallet currency id 31, got %d", ledger.WalletCurrencyID)
				}
				if ledger.Debit != 5_000 || ledger.Credit != 0 {
					t.Fatalf("expected debit 5000 and credit 0, got debit %d credit %d", ledger.Debit, ledger.Credit)
				}
				if ledger.Reference != model.LedgerReferencePayment {
					t.Fatalf("expected reference %s, got %s", model.LedgerReferencePayment, ledger.Reference)
				}
				return &ledgerID, nil
			},
		},
	}

	recorder := httptest.NewRecorder()
	request := requestWithVars(
		http.MethodPost,
		"/wallets/10/pay",
		`{"currency_code":"USD","amount":"50.00"}`,
		map[string]string{"id": "10"},
	)
	handler.PayWithWalletHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
}

func TestTransferWalletHandler(t *testing.T) {
	db, sqlMock, cleanup := newHandlerTestDB(t)
	defer cleanup()

	updatedAt := time.Date(2026, 7, 5, 11, 0, 0, 0, time.UTC)
	ledgerID := int64(501)
	insertedLedgers := make([]model.Ledger, 0, 2)
	sqlMock.ExpectBegin()
	expectWalletWriteLock(sqlMock)
	sqlMock.ExpectCommit()

	handler := &WalletHandler{
		db: db,
		walletRepository: &mockWalletRepository{
			findByIDFunc: func(trx *sqlx.Tx, id int64) (*model.Wallet, error) {
				switch id {
				case 10:
					return &model.Wallet{ID: 10, UserID: 7, Status: model.WalletStatusActivate}, nil
				case 20:
					return &model.Wallet{ID: 20, UserID: 8, Status: model.WalletStatusActivate}, nil
				default:
					t.Fatalf("unexpected wallet id %d", id)
					return nil, nil
				}
			},
		},
		walletCurrencyRepository: &mockWalletCurrencyRepository{
			findByWalletIDAndCurrencyFunc: func(
				trx *sqlx.Tx,
				walletID int64,
				currency model.Currency,
			) (*model.WalletCurrency, error) {
				if currency != model.CurrencyUSD {
					t.Fatalf("expected currency %s, got %s", model.CurrencyUSD, currency)
				}
				switch walletID {
				case 10:
					return &model.WalletCurrency{
						ID:                      31,
						WalletID:                10,
						Currency:                model.CurrencyUSD,
						ClosingBalance:          20_000,
						ClosingBalanceUpdatedAt: updatedAt,
					}, nil
				case 20:
					return &model.WalletCurrency{ID: 41, WalletID: 20, Currency: model.CurrencyUSD}, nil
				default:
					t.Fatalf("unexpected wallet id %d", walletID)
					return nil, nil
				}
			},
		},
		ledgerRepository: &mockLedgerRepository{
			sumOneByWalletCurrencyFunc: func(trx *sqlx.Tx, walletCurrencyID int64, afterCreatedAt time.Time) (*model.LedgerSum, error) {
				if walletCurrencyID != 31 {
					t.Fatalf("expected wallet currency id 31, got %d", walletCurrencyID)
				}
				if !afterCreatedAt.Equal(updatedAt) {
					t.Fatalf("expected after created at %s, got %s", updatedAt, afterCreatedAt)
				}
				return &model.LedgerSum{WalletCurrencyID: 31, Debit: 0, Credit: 0}, nil
			},
			insertFunc: func(trx *sqlx.Tx, ledger *model.Ledger) (*int64, error) {
				insertedLedgers = append(insertedLedgers, *ledger)
				return &ledgerID, nil
			},
		},
	}

	recorder := httptest.NewRecorder()
	request := requestWithVars(
		http.MethodPost,
		"/wallets/10/transfer/USD",
		`{"wallet_id_destination":20,"currency_code_destination":"USD","amount":"75.00"}`,
		map[string]string{"id": "10", "country_code": "USD"},
	)
	handler.TransferWalletHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	if len(insertedLedgers) != 2 {
		t.Fatalf("expected 2 inserted ledgers, got %d", len(insertedLedgers))
	}
	if insertedLedgers[0].WalletCurrencyID != 31 ||
		insertedLedgers[0].Debit != 7_500 ||
		insertedLedgers[0].Credit != 0 ||
		insertedLedgers[0].Reference != model.LedgerReferenceTransfer {
		t.Fatalf("unexpected source ledger: %#v", insertedLedgers[0])
	}
	if insertedLedgers[1].WalletCurrencyID != 41 ||
		insertedLedgers[1].Debit != 0 ||
		insertedLedgers[1].Credit != 7_500 ||
		insertedLedgers[1].Reference != model.LedgerReferenceReceive {
		t.Fatalf("unexpected destination ledger: %#v", insertedLedgers[1])
	}
}

func TestSuspendWalletHandler(t *testing.T) {
	db, sqlMock, cleanup := newHandlerTestDB(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	expectWalletWriteLock(sqlMock)
	sqlMock.ExpectCommit()

	handler := &WalletHandler{
		db: db,
		walletRepository: &mockWalletRepository{
			findByIDFunc: func(trx *sqlx.Tx, id int64) (*model.Wallet, error) {
				if id != 10 {
					t.Fatalf("expected wallet id 10, got %d", id)
				}
				return &model.Wallet{ID: 10, UserID: 7, Status: model.WalletStatusActivate}, nil
			},
			suspendFunc: func(trx *sqlx.Tx, id int64) error {
				if id != 10 {
					t.Fatalf("expected wallet id 10, got %d", id)
				}
				return nil
			},
		},
	}

	recorder := httptest.NewRecorder()
	request := requestWithVars(http.MethodPost, "/wallets/10/suspend", "", map[string]string{"id": "10"})
	handler.SuspendWalletHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
}

func TestUpdateCloseBalanceHandler(t *testing.T) {
	db, sqlMock, cleanup := newHandlerTestDB(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	expectWalletWriteLock(sqlMock)
	sqlMock.ExpectCommit()

	handler := &WalletHandler{
		db: db,
		walletCurrencyRepository: &mockWalletCurrencyRepository{
			updateClosingBalanceFunc: func(trx *sqlx.Tx) error {
				if trx == nil {
					t.Fatal("expected transaction")
				}
				return nil
			},
		},
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/wallet/update_close_balance", nil)
	handler.UpdateCloseBalanceHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
}

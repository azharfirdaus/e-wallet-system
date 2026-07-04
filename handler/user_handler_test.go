package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	handlermodel "github.com/azharfirdaus/e-wallet-system/handler/model"
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/jmoiron/sqlx"
)

type mockUserRepository struct {
	insertFunc func(trx *sqlx.Tx, user *model.User) (*int64, error)
}

func (m *mockUserRepository) Insert(trx *sqlx.Tx, user *model.User) (*int64, error) {
	if m.insertFunc == nil {
		panic("mockUserRepository.Insert is not implemented")
	}
	return m.insertFunc(trx, user)
}

func newHandlerTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}

	return sqlx.NewDb(db, "sqlmock"), mock, func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sql expectations: %v", err)
		}
		_ = db.Close()
	}
}

func TestCreateUserHandler(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		db, sqlMock, cleanup := newHandlerTestDB(t)
		defer cleanup()

		id := int64(42)
		sqlMock.ExpectBegin()
		sqlMock.ExpectCommit()

		handler := &UserHandler{
			db: db,
			userRepository: &mockUserRepository{
				insertFunc: func(trx *sqlx.Tx, user *model.User) (*int64, error) {
					if trx == nil {
						t.Fatal("expected transaction")
					}
					if user == nil {
						t.Fatal("expected user")
					}
					return &id, nil
				},
			},
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/users", nil)
		handler.CreateUserHandler(recorder, request)

		if recorder.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
		}
		if recorder.Header().Get("Content-Type") != "application/json" {
			t.Fatalf("expected content type application/json, got %s", recorder.Header().Get("Content-Type"))
		}

		var response handlermodel.CreateUserResponse
		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if response.ID != id {
			t.Fatalf("expected id %d, got %d", id, response.ID)
		}
	})

	t.Run("repository error rolls back", func(t *testing.T) {
		db, sqlMock, cleanup := newHandlerTestDB(t)
		defer cleanup()

		sqlMock.ExpectBegin()
		sqlMock.ExpectRollback()

		handler := &UserHandler{
			db: db,
			userRepository: &mockUserRepository{
				insertFunc: func(trx *sqlx.Tx, user *model.User) (*int64, error) {
					return nil, errors.New("insert user")
				},
			},
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/users", nil)
		handler.CreateUserHandler(recorder, request)

		if recorder.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
		}
	})
}

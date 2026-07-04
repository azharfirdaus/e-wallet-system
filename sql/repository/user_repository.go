package repository

import (
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Insert(m *model.User) (*int64, error)
}

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (u *UserRepositoryImpl) Insert(m *model.User) (*int64, error) {
	trx, err := u.db.Beginx()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO public.users DEFAULT VALUES
		RETURNING id
	`

	var id int64
	if err := trx.QueryRowx(query).Scan(&id); err != nil {
		_ = trx.Rollback()
		return nil, err
	}

	if err := trx.Commit(); err != nil {
		return nil, err
	}

	return &id, nil
}

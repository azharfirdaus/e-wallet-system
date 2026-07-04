package repository

import (
	"github.com/azharfirdaus/e-wallet-system/sql/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Insert(trx *sqlx.Tx, m *model.User) (*int64, error)
}

type UserRepositoryImpl struct{}

func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (u *UserRepositoryImpl) Insert(trx *sqlx.Tx, m *model.User) (*int64, error) {
	query := `
		INSERT INTO public.users DEFAULT VALUES
		RETURNING id
	`

	var id int64
	if err := trx.QueryRowx(query).Scan(&id); err != nil {
		return nil, err
	}

	return &id, nil
}

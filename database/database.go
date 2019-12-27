package database

import (
	users "github.com/microapis/users-api"
	"github.com/microapis/users-api/database/postgres"
	"github.com/jmoiron/sqlx"
)

// Store ...
type Store interface {
	UserGet(*users.Query) (*users.User, error)
	UserCreate(*users.User) error
	UserList() ([]*users.User, error)

	// TODO(ca): below methods are not implemented
	// Update(*users.User) error
	// Delete(*users.User) error
}

// NewPostgres ...
func NewPostgres(dsn string) (Store, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &postgres.UserStore{
		Store: db,
	}, nil
}

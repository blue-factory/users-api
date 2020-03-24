package database

import (
	"github.com/jmoiron/sqlx"
	users "github.com/microapis/users-api"
	"github.com/microapis/users-api/database/postgres"
)

// Store ...
type Store interface {
	UserGet(*users.Query) (*users.User, error)
	UserCreate(*users.User) error
	UserList() ([]*users.User, error)
	Update(*users.User) error
	Delete(*users.User) error
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

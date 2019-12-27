package service

import (
	"errors"

	users "github.com/microapis/users-api"
	"github.com/microapis/users-api/database"
)

// NewUsers ...
func NewUsers(store database.Store) *Users {
	return &Users{
		Store: store,
	}
}

// Users ...
type Users struct {
	Store database.Store
}

// GetByID ...
func (us *Users) GetByID(id string) (*users.User, error) {
	q := &users.Query{
		ID: id,
	}

	return us.Store.UserGet(q)
}

// GetByEmail ...
func (us *Users) GetByEmail(email string) (*users.User, error) {
	q := &users.Query{
		Email: email,
	}

	return us.Store.UserGet(q)
}

// Create ...
func (us *Users) Create(u *users.User) error {
	return us.Store.UserCreate(u)
}

// Update ...
func (us *Users) Update(*users.User) error {
	return errors.New("methods is not implemented")
}

// Delete ...
func (us *Users) Delete(*users.User) error {
	return errors.New("method is not implemented")
}

// UserList ...
func (us *Users) UserList() ([]*users.User, error) {
	return us.Store.UserList()
}

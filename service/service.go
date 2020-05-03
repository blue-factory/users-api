package service

import (
	users "github.com/microapis/users-api"
	"github.com/microapis/users-api/database"
)

// NewUsers ...
func NewUsers(store database.Store, events *users.Events) *Users {
	return &Users{
		Store:  store,
		Events: events,
	}
}

// Users ...
type Users struct {
	Store  database.Store
	Events *users.Events
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
	// before create user event
	if us.Events.BeforeCreate != nil {
		err := us.Events.BeforeCreate()
		if err != nil {
			return err
		}
	}

	err := us.Store.UserCreate(u)
	if err != nil {
		return err
	}

	// after create user event
	if us.Events.AfterCreate != nil {
		err := us.Events.AfterCreate()
		if err != nil {
			return err
		}
	}

	return nil
}

// Update ...
func (us *Users) Update(u *users.User) error {
	return us.Store.Update(u)
}

// Delete ...
func (us *Users) Delete(u *users.User) error {
	return us.Store.Delete(u)
}

// UserList ...
func (us *Users) UserList() ([]*users.User, error) {
	return us.Store.UserList()
}

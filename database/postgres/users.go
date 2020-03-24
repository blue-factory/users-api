package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	users "github.com/microapis/users-api"
)

// UserStore ...
type UserStore struct {
	Store *sqlx.DB
}

// UserGet ...
func (us *UserStore) UserGet(q *users.Query) (*users.User, error) {
	query := squirrel.Select("*").From("users").Where("deleted_at is null")

	if q.ID == "" && q.Email == "" {
		return nil, errors.New("must proovide a query")
	}

	if q.ID != "" {
		query = query.Where("id = ?", q.ID)
	}

	if q.Email != "" {
		query = query.Where("email = ?", q.Email)
	}

	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := us.Store.QueryRowx(sql, args...)

	c := &users.User{}
	if err := row.StructScan(c); err != nil {
		return nil, err
	}

	return c, nil
}

// UserCreate ...
func (us *UserStore) UserCreate(u *users.User) error {

	sql, args, err := squirrel.
		Insert("users").
		Columns("email", "name", "password").
		Values(u.Email, u.Name, u.Password).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	row := us.Store.QueryRowx(sql, args...)
	if err := row.StructScan(u); err != nil {
		return err
	}

	return nil
}

// UserList ...
func (us *UserStore) UserList() ([]*users.User, error) {
	query := squirrel.Select("*").From("users").Where("deleted_at is null")

	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := us.Store.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	uu := make([]*users.User, 0)

	for rows.Next() {
		u := &users.User{}
		if err := rows.StructScan(u); err != nil {
			return nil, err
		}

		uu = append(uu, u)
	}

	return uu, nil
}

// Update ...
func (us *UserStore) Update(u *users.User) error {
	sql, args, err := squirrel.Update("users").Set("email", u.Email).Set("name", u.Name).Set("password", u.Password).Where("id = ?", u.ID).Suffix("returning *").PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return err
	}

	row := us.Store.QueryRowx(sql, args...)
	return row.StructScan(u)
}

// Delete ...
func (us *UserStore) Delete(u *users.User) error {
	row := us.Store.QueryRowx("update users set deleted_at = $1 where id = $2 returning *", time.Now(), u.ID)

	if err := row.StructScan(u); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	return nil
}

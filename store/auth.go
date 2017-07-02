package store

import (
	"database/sql"
	"errors"

	"github.com/kusubooru/shimmie"
)

// Errors returned by Verify.
var (
	ErrWrongCredentials = errors.New("wrong username or password")
	ErrNotFound         = errors.New("entry not found")
)

// Verify compares the provided username and password with the username and
// password hash stored in the shimmie database.
//
// It can return:
//
// - The shimmie User on success.
//
// - ErrNotFound if the username does not exist.
//
// - ErrWrongCredentials if the username and password do not match.
//
// - An error if something goes wrong with the database.
func (db *datastore) Verify(username, password string) (*shimmie.User, error) {
	u, err := db.GetUserByName(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if u.Pass == shimmie.PasswordHash(username, password) {
		return u, nil
	}
	return nil, ErrWrongCredentials
}

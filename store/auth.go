package store

import (
	"database/sql"

	"github.com/kusubooru/shimmie"
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
			return nil, shimmie.ErrNotFound
		}
		return nil, err
	}
	if u.Pass == shimmie.PasswordHash(username, password) {
		return u, nil
	}
	return nil, shimmie.ErrWrongCredentials
}

package store

import (
	"database/sql"

	"github.com/kusubooru/shimmie"
)

func (db *datastore) GetUser(username string) (*shimmie.User, error) {
	var (
		u     shimmie.User
		pass  sql.NullString
		email sql.NullString
	)
	err := db.QueryRow(userGetQuery, username).Scan(
		&u.ID,
		&u.Name,
		&pass,
		&u.JoinDate,
		&u.Admin,
		&email,
		&u.Class,
	)
	if err != nil {
		return nil, err
	}
	if pass.Valid {
		u.Pass = pass.String
	}
	if email.Valid {
		u.Email = email.String
	}
	return &u, nil
}

func (db *datastore) DeleteUser(id int64) error {
	stmt, err := db.Prepare(userDeleteStmt)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(id); err != nil {
		return err
	}
	return nil
}

func (db *datastore) CreateUser(u *shimmie.User) error {
	stmt, err := db.Prepare(userInsertStmt)
	if err != nil {
		return err
	}
	if u.Class == "" {
		u.Class = "user"
	}
	if u.Admin == "" {
		u.Admin = "N"
	}
	hash := shimmie.Hash(u.Pass)
	_, err = stmt.Exec(u.Name, hash, u.Email, u.Class)
	if err != nil {
		return err
	}

	// Executing an extra statement in order to get the value of JoinDate as
	// its value is created by the database with NOW(). Creating the value of
	// JoinDate with time.Now() beforehand, results in slightly different dates
	// as time.Time has 9 point decimal precision while MySQL DATETIME has 6 at
	// max.
	//
	// If a better solution for getting the JoinDate value can be found then we
	// can avoid the extra execution by simply filling the value of ID:
	//
	//    id, err := res.LastInsertId()
	//    if err != nil {
	//		return err
	//    }
	//    u.ID = id
	storedUser, err := db.GetUser(u.Name)
	if err != nil {
		return err
	}
	*u = *storedUser

	return nil
}

const (
	userGetQuery = `
SELECT * 
FROM users 
WHERE name = ?
`
	userInsertStmt = `
INSERT users
SET
  name=?,
  pass=?,
  joindate=NOW(),
  class=?,
  email=?
`
	userDeleteStmt = `
DELETE
FROM users
WHERE id = ?
`
)

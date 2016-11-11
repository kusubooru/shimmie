package store

import (
	"database/sql"

	"github.com/kusubooru/shimmie"
)

func (db *datastore) GetUser(userID int64) (*shimmie.User, error) {
	return db.getUserBy(userGetQuery, userID)
}

func (db *datastore) GetUserByName(username string) (*shimmie.User, error) {
	return db.getUserBy(userGetByNameQuery, username)
}

func (db *datastore) getUserBy(query string, id interface{}) (*shimmie.User, error) {
	var (
		u     shimmie.User
		pass  sql.NullString
		email sql.NullString
	)
	err := db.QueryRow(query, id).Scan(
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
	storedUser, err := db.GetUserByName(u.Name)
	if err != nil {
		return err
	}
	*u = *storedUser

	return nil
}

func (db *datastore) CountUsers() (int, error) {
	return count(db.DB, userCountQuery)
}

func count(db *sql.DB, query string) (int, error) {
	var (
		count int
	)
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *datastore) GetAllUsers(limit, offset int) ([]shimmie.User, error) {
	if limit < 0 {
		count, cerr := db.CountUsers()
		if cerr != nil {
			return nil, cerr
		}
		limit = count
	}
	rows, err := db.Query(userGetAllQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); err == nil {
			err = cerr
			return
		}
	}()

	var (
		u     shimmie.User
		users []shimmie.User
		pass  sql.NullString
		email sql.NullString
	)
	for rows.Next() {
		err = rows.Scan(
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
		users = append(users, u)
	}
	return users, err
}

const (
	userGetQuery = `
SELECT *
FROM users
WHERE id = ?
`
	userGetByNameQuery = `
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

	userCountQuery = `
SELECT COUNT(*)
FROM users
`

	userGetAllQuery = `
SELECT *
FROM users
LIMIT ? OFFSET ?
`
)

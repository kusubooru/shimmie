package shimmiedb

import (
	"database/sql"
	"fmt"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

// Open creates a database connection for the given driver and configuration.
func Open(dataSource string, pingRetries int) (*DB, error) {
	db, err := openDB(dataSource, pingRetries)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// openDB opens a new database connection with the specified connection string.
func openDB(dataSource string, pingRetries int) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, fmt.Errorf("opening db: %v", err)
	}

	// per issue https://github.com/go-sql-driver/mysql/issues/257
	db.SetMaxIdleConns(0)

	if err := pingDatabase(db, pingRetries); err != nil {
		return nil, fmt.Errorf("pinging db %d times: %v", pingRetries, err)
	}
	return db, nil
}

// helper function to ping the database with backoff to ensure a connection can
// be established before we proceed.
func pingDatabase(db *sql.DB, pingRetries int) (err error) {
	for i := 0; i < pingRetries; i++ {
		err = db.Ping()
		if err == nil {
			return
		}
		//log.Printf("database ping failed, retry in 1s: %v", err)
		time.Sleep(time.Second)
	}
	return
}

func (db DB) Close() error {
	return db.DB.Close()
}

// Tx allows to perform a function in a transaction. It detects error and panic
// and in that case it rollbacks otherwise it commits the transaction.
func Tx(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	return txFunc(tx)
}

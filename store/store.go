package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type Datastore struct {
	*sql.DB
}

// Open creates a database connection for the given driver and configuration.
func Open(config string, pingRetries int) *Datastore {
	db := openDB(config, pingRetries)
	return &Datastore{db}
}

// openDB opens a new database connection with the specified driver and
// connection string.
func openDB(config string, pingRetries int) *sql.DB {
	db, err := sql.Open("mysql", config)
	if err != nil {
		log.Fatalln("database connection failed:", err)
	}

	// per issue https://github.com/go-sql-driver/mysql/issues/257
	db.SetMaxIdleConns(0)

	if err := pingDatabase(db, pingRetries); err != nil {
		log.Fatalln("database ping attempts failed:", err)
	}
	return db
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

func (db Datastore) Close() error {
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

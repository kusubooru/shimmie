package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/kusubooru/shimmie"
)

type datastore struct {
	*sql.DB
}

func (db datastore) SQLDB() *sql.DB {
	return db.DB
}

// Open creates a database connection for the given driver and configuration.
func Open(driver, config string) shimmie.Store {
	db := openDB(driver, config)
	return &datastore{db}
}

// openDB opens a new database connection with the specified driver and
// connection string.
func openDB(driver, config string) *sql.DB {
	db, err := sql.Open(driver, config)
	if err != nil {
		log.Fatalln("database connection failed:", err)
	}
	if driver == "mysql" {
		// per issue https://github.com/go-sql-driver/mysql/issues/257
		db.SetMaxIdleConns(0)
	}
	if err := pingDatabase(db); err != nil {
		log.Fatalln("database ping attempts failed:", err)
	}
	return db
}

// helper function to ping the database with backoff to ensure a connection can
// be established before we proceed.
func pingDatabase(db *sql.DB) (err error) {
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			return
		}
		log.Print("database ping failed. retry in 1s")
		time.Sleep(time.Second)
	}
	return
}

func (db datastore) Close() error {
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

package store_test

import (
	"database/sql"
	"flag"
	"fmt"
	"testing"

	"github.com/kusubooru/shimmie"
	"github.com/kusubooru/shimmie/store"
)

var (
	dbDriver     = flag.String("dbdriver", "mysql", "database driver")
	dbInitConfig = flag.String("dbinitconfig", "", "username:password@(host:port)/database?parseTime=true")
	dbConfig     = flag.String("dbconfig", "", "username:password@(host:port)/database?parseTime=true")
	dbName       = flag.String("dbname", "kusubooru_test_db", "database name to run tests on")
)

func init() {
	flag.Parse()
	if *dbInitConfig == "" {
		*dbInitConfig = "kusubooru:kusubooru@/?parseTime=true"
	}
	if *dbConfig == "" {
		*dbConfig = fmt.Sprintf("kusubooru:kusubooru@/%s?parseTime=true", *dbName)
	}
}

func setup(t *testing.T) shimmie.Store {
	create(t)
	shim := store.Open(*dbDriver, *dbConfig)
	insert(t, shim.SQLDB())
	return shim
}

func insert(t *testing.T, db *sql.DB) {
	//tx, err := db.Begin()
	//if err != nil {
	//	t.Fatalf("could not start insert transaction: %s", err)
	//}
	//defer func() {
	//	if rerr := tx.Rollback(); rerr != nil {
	//		t.Errorf("could not rollback insert transaction: %s", rerr)
	//	}
	//}()
	//_, err = tx.Exec("insert into ...")
	//if err != nil {
	//	t.Errorf("could not insert: %s", err)
	//}
}

func create(t *testing.T) {
	shim := store.Open(*dbDriver, *dbInitConfig)
	db := shim.SQLDB()
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("could not close db connection: %s", err)
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("could not start transaction: %s", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(fmt.Sprintf("CREATE DATABASE %s;", *dbName))
	if err != nil {
		t.Errorf("could not create test db: %s", err)
	}
	_, err = tx.Exec(fmt.Sprintf("USE %s;", *dbName))
	if err != nil {
		t.Errorf("could not select test db: %s", err)
	}
	loadSchema(t, tx)
	if err := tx.Commit(); err != nil {
		t.Errorf("could not commit transaction: %s", err)
	}
}

func loadSchema(t *testing.T, tx *sql.Tx) {
	_, err := tx.Exec(usersCreateQuery)
	if err != nil {
		t.Errorf("could not create table users: %s", err)
	}

	_, err = tx.Exec(imagesCreateQuery)
	if err != nil {
		t.Errorf("could not create table images: %s", err)
	}

	_, err = tx.Exec(tagHistoriesCreateQuery)
	if err != nil {
		t.Errorf("could not create table tag_histories: %s", err)
	}
}

func teardown(t *testing.T, db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf("DROP DATABASE %s;", *dbName))
	if err != nil {
		t.Fatalf("could not drop test db: %s", err)
	}
}

const (
	// datetime must have default value
	usersCreateQuery = `
create table users (
	id INTEGER NOT NULL AUTO_INCREMENT,
	name VARCHAR(32) UNIQUE NOT NULL,
	pass VARCHAR(250),
	joindate datetime NOT NULL ,
	class VARCHAR(32) NOT NULL DEFAULT 'user',
	email VARCHAR(128),
	PRIMARY KEY (id)
)
	`
	tagHistoriesCreateQuery = `
create table tag_histories (
	id INTEGER NOT NULL AUTO_INCREMENT,
	image_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	user_ip CHAR(15) NOT NULL,
	tags TEXT NOT NULL,
	date_set DATETIME NOT NULL,
	PRIMARY KEY (id),
	FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
 `
	// datetime default
	imagesCreateQuery = `
create table images (
	id INTEGER NOT NULL AUTO_INCREMENT,
	owner_id INTEGER NOT NULL,
	owner_ip char(15) NOT NULL,
	filename VARCHAR(64) NOT NULL,
	filesize INTEGER NOT NULL,
	hash CHAR(32) UNIQUE NOT NULL,
	ext CHAR(4) NOT NULL,
	source VARCHAR(255),
	width INTEGER NOT NULL,
	height INTEGER NOT NULL,
	posted DATETIME NOT NULL ,
	locked enum('Y','N') NOT NULL DEFAULT 'N',
	PRIMARY KEY (id),
	FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE RESTRICT
)
 `
)

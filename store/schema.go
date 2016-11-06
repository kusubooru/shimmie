package store

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/kusubooru/shimmie"
)

type schema struct {
	*sql.DB
}

// NewSchemer returns an implementation of Schemer that allows to easily create
// and drop the database schema.
func NewSchemer(driverName, username, password, host, port string) shimmie.Schemer {
	db := connect(driverName, username, password, host, port)
	return &schema{db}
}

func connect(driverName, username, password, host, port string) *sql.DB {
	dataSourceName := fmt.Sprintf("%s:%s@(%s:%s)/?parseTime=true", username, password, host, port)
	return openDB(driverName, dataSourceName)
}

func (db schema) Create(dbName string) error {
	return Tx(db.DB, func(tx *sql.Tx) error {
		if _, err := tx.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName)); err != nil {
			return fmt.Errorf("could not create db: %s", err)
		}

		if _, err := tx.Exec(fmt.Sprintf("USE %s;", dbName)); err != nil {
			return fmt.Errorf("could not select db: %s", err)
		}

		if query, err := createSchema(tx); err != nil {
			return fmt.Errorf("failed to execute query:\n\n  %q\n\n  Reason:\n\n  %v\n", query, err)
		}
		return nil
	})
}

func (db schema) Drop(dbName string) error {
	return Tx(db.DB, func(tx *sql.Tx) error {
		if _, err := tx.Exec(fmt.Sprintf("DROP DATABASE %s;", dbName)); err != nil {
			return fmt.Errorf("could not drop db %s: %v", dbName, err)
		}
		return nil
	})
}

func (db schema) Close() error {
	return db.DB.Close()
}

// MySQL specific error for when we try to run alter queries to add new a
// column and the column already exists.
//
// See: https://github.com/VividCortex/mysqlerr/blob/master/mysqlerr.go
const duplicateColumnName = 1060

func createSchema(tx *sql.Tx) (string, error) {
	for _, query := range createStatements {
		if _, err := tx.Exec(query); err != nil {
			return query, err
		}
	}

	for _, query := range alterStatements {
		if _, err := tx.Exec(query); err != nil {
			if driverErr, ok := err.(*mysql.MySQLError); ok {
				if driverErr.Number != duplicateColumnName {
					return query, err
				}
			}
		}
	}
	return "", nil
}

var createStatements = []string{
	usersCreateTableStmt,
	imagesCreateTableStmt,
	tagsCreateTableStmt,
	tagHistoriesCreateTableStmt,
	imageTagsCreateTableStmt,
	aliasesCreateTableStmt,
}

var alterStatements = []string{}

const (
	usersCreateTableStmt = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL AUTO_INCREMENT,
	name VARCHAR(32) NOT NULL,
	pass VARCHAR(250),
	joindate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	admin enum('Y','N') NOT NULL DEFAULT 'N',
	class VARCHAR(32) NOT NULL DEFAULT 'user',
	email VARCHAR(128),
	PRIMARY KEY (id),
	UNIQUE KEY (name)
);
`
	imagesCreateTableStmt = `
CREATE TABLE IF NOT EXISTS images (
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
	posted DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	locked enum('Y','N') NOT NULL DEFAULT 'N',
	PRIMARY KEY (id),
	FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE RESTRICT
);
`
	tagsCreateTableStmt = `
CREATE TABLE IF NOT EXISTS tags (
  id int(11) NOT NULL AUTO_INCREMENT,
  tag varchar(64) NOT NULL,
  ` + "`count`" + `int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (id),
  UNIQUE KEY tag (tag)
);
`
	tagHistoriesCreateTableStmt = `
CREATE TABLE IF NOT EXISTS tag_histories (
	id INTEGER NOT NULL AUTO_INCREMENT,
	image_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	user_ip CHAR(15) NOT NULL,
	tags TEXT NOT NULL,
	date_set DATETIME NOT NULL,
	PRIMARY KEY (id),
	FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
`
	imageTagsCreateTableStmt = `
CREATE TABLE IF NOT EXISTS image_tags (
  image_id int(11) NOT NULL,
  tag_id int(11) NOT NULL,
  UNIQUE KEY image_tags_id (image_id,tag_id),
  INDEX image_id (image_id),
  INDEX tag_id (tag_id),
  CONSTRAINT image_tags_ibfk_1 FOREIGN KEY (image_id) REFERENCES images (id) ON DELETE CASCADE,
  CONSTRAINT image_tags_ibfk_2 FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
);
`
	aliasesCreateTableStmt = `
CREATE TABLE IF NOT EXISTS aliases (
  oldtag varchar(128) NOT NULL,
  newtag varchar(128) NOT NULL,
  PRIMARY KEY (oldtag),
  INDEX newtag (newtag)
);
`
)

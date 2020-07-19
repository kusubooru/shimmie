package shimmiedb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// Schema has methods which allow to create the db schema as well as truncate
// all tables.
type Schema struct {
	*sql.DB
}

// NewSchemer returns an implementation of Schemer that allows to easily create
// and drop the database schema.
func NewSchemer(dataSource string, pingRetries int) (*Schema, error) {
	db, err := openDB(dataSource, pingRetries)
	if err != nil {
		return nil, err
	}
	return &Schema{db}, nil
}

// Create creates the database schema.
func (db Schema) Create() error {
	return Tx(db.DB, func(tx *sql.Tx) error {
		if query, err := createSchema(tx); err != nil {
			return fmt.Errorf("failed to execute query:\n%s\nReason: %v", query, err)
		}
		return nil
	})
}

func (db Schema) allTables(ctx context.Context) ([]string, error) {
	const q = `
	SELECT table_name
	FROM information_schema.tables
	WHERE table_schema="shimmie";`

	rows, err := db.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss := []string{}
	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}
	return ss, rows.Err()

}

// TruncateTables truncates all the database tables.
func (db Schema) TruncateTables(ctx context.Context) error {
	tables, err := db.allTables(ctx)
	if err != nil {
		return fmt.Errorf("fetching all tables: %v", err)
	}

	b := strings.Builder{}
	b.WriteString("SET FOREIGN_KEY_CHECKS=0;\n")
	for _, t := range tables {
		b.WriteString(fmt.Sprintf("TRUNCATE TABLE %s;\n", t))
	}
	b.WriteString("SET FOREIGN_KEY_CHECKS=1;")

	query := b.String()
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("truncating all tables using query:\n%s\nResult: %v", query, err)
	}
	return nil
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
	privateMessageTableStmt,
}

var alterStatements = []string{}

const (
	usersCreateTableStmt = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL AUTO_INCREMENT,
	name VARCHAR(32) NOT NULL,
	pass VARCHAR(250),
	joindate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
	posted TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
	privateMessageTableStmt = `
CREATE TABLE IF NOT EXISTS private_message (
  id int(11) NOT NULL AUTO_INCREMENT,
  from_id int(11) NOT NULL,
  from_ip char(15) NOT NULL,
  to_id int(11) NOT NULL,
  sent_date datetime NOT NULL,
  subject varchar(64) NOT NULL,
  message text NOT NULL,
  is_read enum('Y','N') NOT NULL DEFAULT 'N',
  PRIMARY KEY (id),
  KEY to_id (to_id),
  KEY from_id (from_id),
  CONSTRAINT private_message_ibfk_1 FOREIGN KEY (from_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT private_message_ibfk_2 FOREIGN KEY (to_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`
)

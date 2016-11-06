package store

import "github.com/kusubooru/shimmie"

func (db datastore) GetAlias(oldTag string) (*shimmie.Alias, error) {
	var (
		a shimmie.Alias
	)
	err := db.QueryRow(aliasGetQuery, oldTag).Scan(
		&a.OldTag,
		&a.NewTag,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (db datastore) DeleteAlias(oldTag string) error {
	stmt, err := db.Prepare(aliasDeleteStmt)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(oldTag); err != nil {
		return err
	}
	return nil
}

func (db datastore) CreateAlias(alias *shimmie.Alias) error {
	stmt, err := db.Prepare(aliasInsertStmt)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(alias.NewTag, alias.OldTag)
	if err != nil {
		return err
	}
	return nil
}

const (
	aliasGetQuery = `
SELECT *
FROM aliases
WHERE oldtag = ?
`
	aliasDeleteStmt = `
DELETE
FROM aliases
WHERE oldtag = ?
`

	aliasInsertStmt = `
INSERT aliases
SET
  newtag = ?,
  oldtag = ?
`
)

package shimmiedb

import "github.com/kusubooru/shimmie"

// GetAlias returns an alias based on its old tag.
func (db *DB) GetAlias(oldTag string) (*shimmie.Alias, error) {
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

// DeleteAlias deletes an alias based on its old tag.
func (db *DB) DeleteAlias(oldTag string) error {
	stmt, err := db.Prepare(aliasDeleteStmt)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(oldTag); err != nil {
		return err
	}
	return nil
}

// CreateAlias creates a new alias.
func (db *DB) CreateAlias(alias *shimmie.Alias) error {
	stmt, err := db.Prepare(aliasInsertStmt)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(alias.NewTag, alias.OldTag)
	return err
}

// CountAlias returns how many alias entries exist in the database.
func (db *DB) CountAlias() (int, error) {
	var (
		count int
	)
	err := db.QueryRow(aliasCountQuery).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// GetAllAlias returns alias entries of the database based on a limit and
// an offset. If limit < 0, CountAlias will also be executed to get the
// maximum limit and return all alias entries. Offset still works in this
// case. For example, assuming 10 entries, GetAllAlias(-1, 0), will return
// all 10 entries and GetAllAlias(-1, 8) will return the last 2 entries.
func (db *DB) GetAllAlias(limit, offset int) ([]shimmie.Alias, error) {
	if limit < 0 {
		count, cerr := db.CountAlias()
		if cerr != nil {
			return nil, cerr
		}
		limit = count
	}
	rows, err := db.Query(aliasGetAllQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); err == nil {
			err = cerr
			return
		}
	}()

	var alias []shimmie.Alias
	for rows.Next() {
		var a shimmie.Alias
		err = rows.Scan(
			&a.OldTag,
			&a.NewTag,
		)
		if err != nil {
			return nil, err
		}
		alias = append(alias, a)
	}
	return alias, err
}

// FindAlias returns all alias matching an oldTag or a newTag or both.
func (db *DB) FindAlias(oldTag, newTag string) ([]shimmie.Alias, error) {
	rows, err := db.Query(aliasFindQuery, "%"+oldTag+"%", "%"+newTag+"%")
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); err == nil {
			err = cerr
			return
		}
	}()

	var alias []shimmie.Alias
	for rows.Next() {
		var a shimmie.Alias
		err = rows.Scan(
			&a.OldTag,
			&a.NewTag,
		)
		if err != nil {
			return nil, err
		}
		alias = append(alias, a)
	}
	return alias, err
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

	aliasCountQuery = `
SELECT COUNT(*)
FROM aliases
`

	aliasGetAllQuery = `
SELECT *
FROM aliases
LIMIT ? OFFSET ?
`

	aliasFindQuery = `
SELECT *
FROM aliases
WHERE oldtag like ?
  AND newtag like ?
`
)

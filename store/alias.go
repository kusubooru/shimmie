package store

import "github.com/kusubooru/shimmie"

func (db *datastore) GetAlias(oldTag string) (*shimmie.Alias, error) {
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

func (db *datastore) DeleteAlias(oldTag string) error {
	stmt, err := db.Prepare(aliasDeleteStmt)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(oldTag); err != nil {
		return err
	}
	return nil
}

func (db *datastore) CreateAlias(alias *shimmie.Alias) error {
	stmt, err := db.Prepare(aliasInsertStmt)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(alias.NewTag, alias.OldTag)
	return err
}

func (db *datastore) CountAlias() (int, error) {
	var (
		count int
	)
	err := db.QueryRow(aliasCountQuery).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *datastore) GetAllAlias(limit, offset int) ([]shimmie.Alias, error) {
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

	var (
		a     shimmie.Alias
		alias []shimmie.Alias
	)
	for rows.Next() {
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

func (db *datastore) FindAlias(oldTag, newTag string) ([]shimmie.Alias, error) {
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

	var (
		a     shimmie.Alias
		alias []shimmie.Alias
	)
	for rows.Next() {
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

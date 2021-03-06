package shimmiedb

import (
	"github.com/kusubooru/shimmie"
)

// GetTag returns an image tag.
func (db *DB) GetTag(oldTag string) (*shimmie.Tag, error) {
	var (
		t shimmie.Tag
	)
	err := db.QueryRow(tagGetQuery, oldTag).Scan(
		&t.ID,
		&t.Tag,
		&t.Count,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// DeleteTag deletes an image tag.
func (db *DB) DeleteTag(name string) error {
	stmt, err := db.Prepare(tagDeleteStmt)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := stmt.Close(); err == nil {
			err = cerr
			return
		}
	}()
	if _, err := stmt.Exec(name); err != nil {
		return err
	}
	return nil
}

// CreateTag inserts an image tag to the db.
func (db *DB) CreateTag(t *shimmie.Tag) error {
	stmt, err := db.Prepare(tagInsertStmt)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := stmt.Close(); err == nil {
			err = cerr
			return
		}
	}()
	_, err = stmt.Exec(t.Tag, t.Count)
	return err
}

// GetAllTags returns all stored tags.
func (db *DB) GetAllTags(limit, offset int) ([]*shimmie.Tag, error) {
	rows, err := db.Query(tagsGetAllQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); err == nil {
			err = cerr
			return
		}
	}()

	var tags []*shimmie.Tag
	for rows.Next() {
		var t shimmie.Tag
		err = rows.Scan(
			&t.ID,
			&t.Tag,
			&t.Count,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}
	return tags, err
}

const (
	tagGetQuery = `
SELECT *
FROM tags
WHERE tag = ?
`

	tagDeleteStmt = `
DELETE
FROM tags
WHERE tag = ?
`

	tagInsertStmt = `
INSERT tags
SET
  tag = ?,
  count = ?
`
	tagsGetAllQuery = `
SELECT *
FROM tags
LIMIT ? OFFSET ?
`
)

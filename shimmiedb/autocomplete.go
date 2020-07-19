package shimmiedb

import "github.com/kusubooru/shimmie"

func (db *DB) Autocomplete(q string, limit, offset int) ([]*shimmie.Autocomplete, error) {
	if q == "" {
		return []*shimmie.Autocomplete{}, nil
	}
	q = "%" + q + "%"
	rows, err := db.Query(autocompleteQuery, q, q, q, q, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); err == nil {
			err = cerr
			return
		}
	}()

	var autocomplete []*shimmie.Autocomplete
	for rows.Next() {
		var a shimmie.Autocomplete
		err = rows.Scan(&a.Old, &a.Name, &a.Count)
		if err != nil {
			return nil, err
		}
		autocomplete = append(autocomplete, &a)
	}
	return autocomplete, err
}

const autocompleteQuery = `
SELECT a.oldtag AS old, a.newtag AS name, t.count AS count
FROM aliases a, tags t
WHERE
	(a.oldtag LIKE ? AND a.newtag = t.tag) OR
	(t.tag LIKE ? AND a.newtag = t.tag)
UNION
	(SELECT '' AS old, t.tag AS name, t.count AS count
	FROM tags t
	WHERE t.tag LIKE ? AND t.count > 0 AND t.tag NOT IN (
			SELECT oldtag from aliases where oldtag like ?
			UNION
            SELECT newtag from aliases where newtag like ?
		)
    )
ORDER BY count DESC
LIMIT ? OFFSET ?
`

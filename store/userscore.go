package store

import "github.com/kusubooru/shimmie"

func (db *Datastore) MostImageUploads(limit int) ([]shimmie.UserScore, error) {
	const query = `
	SELECT
		count(img.owner_id) as score,
		u.id,
		u.name,
		u.join_date,
		u.email,
		u.class
	FROM images img
	  JOIN users u
	  ON img.owner_id=u.id
	GROUP BY img.owner_id
	ORDER BY count DESC
	LIMIT ?;`

	return db.userScore(query, limit)
}

func (db *Datastore) MostTagChanges(limit int) ([]shimmie.UserScore, error) {
	const query = `
	SELECT
		count(th.user_id) as score,
		u.id,
		u.name,
		u.join_date,
		u.email,
		u.class
	FROM tag_histories th
	  JOIN users u
	  ON th.user_id=u.id
	GROUP BY th.user_id
	ORDER BY count
	LIMIT ?;`

	return db.userScore(query, limit)
}

func (db *Datastore) userScore(query string, limit int) ([]shimmie.UserScore, error) {
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); err == nil {
			err = cerr
			return
		}
	}()

	ss := []shimmie.UserScore{}
	for rows.Next() {
		s := new(shimmie.UserScore)
		err = rows.Scan(
			s.Score,
			s.ID,
			s.Name,
			s.JoinDate,
			s.Email,
			s.Class,
		)
		if err != nil {
			return nil, err
		}
		ss = append(ss, *s)
	}
	err = rows.Err()
	return ss, err
}

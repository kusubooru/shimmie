package store

import "github.com/kusubooru/shimmie"

func (db *datastore) GetTagHistory(imageID int) ([]shimmie.TagHistory, error) {
	rows, err := db.Query(tagHistoryGetQuery, imageID)
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
		th  shimmie.TagHistory
		ths []shimmie.TagHistory
	)
	for rows.Next() {
		err := rows.Scan(
			&th.ID,
			&th.ImageID,
			&th.UserID,
			&th.UserIP,
			&th.Tags,
			&th.DateSet,
			&th.Name,
		)
		if err != nil {
			return nil, err
		}
		ths = append(ths, th)
	}
	return ths, nil
}

const (
	tagHistoryGetQuery = `
SELECT tag_histories.*, users.name
FROM tag_histories
JOIN users ON tag_histories.user_id = users.id
WHERE image_id = ?
ORDER BY tag_histories.id DESC
`
)

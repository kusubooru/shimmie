package store

import (
	"database/sql"

	"github.com/kusubooru/shimmie"
)

func (db *datastore) GetSafeBustedImages(username string) ([]shimmie.Image, error) {
	rows, err := db.Query(imageGetSafeBustedQuery, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		img    shimmie.Image
		images []shimmie.Image
	)
	for rows.Next() {
		var source sql.NullString
		err := rows.Scan(
			&img.ID,
			&img.OwnerID,
			&img.OwnerIP,
			&img.Filename,
			&img.Filesize,
			&img.Hash,
			&img.Ext,
			&source,
			&img.Width,
			&img.Height,
			&img.Posted,
			&img.Locked,
			&img.NumericScore,
			&img.Rating,
			&img.Favorites,
		)
		if err != nil {
			return nil, err
		}
		if source.Valid {
			img.Source = source.String
		}
		images = append(images, img)
	}
	return images, nil
}

// imageGetSafeBustedQuery searches score_log for images set as "Safe" ignoring
// ones from a specific username, extracts the ID of the images and returns
// them.
//
// Warning: MySQL specific query.
const imageGetSafeBustedQuery = `
SELECT
  img.*
FROM images as img, (
  SELECT SUBSTRING_INDEX(SUBSTRING_INDEX(message, '#', -1), ' ', 1) AS id
  FROM score_log
  WHERE message
  LIKE "%set to: Safe"
  AND username != ?
  ORDER BY date_sent DESC) as safe
WHERE img.id = safe.id
AND rating = 's'
`

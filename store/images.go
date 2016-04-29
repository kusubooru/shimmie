package store

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kusubooru/shimmie"
)

func (db *datastore) GetImage(id int) (*shimmie.Image, error) {
	var (
		img    shimmie.Image
		source sql.NullString
	)
	err := db.QueryRow(imageGetQuery, id).Scan(
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

	return &img, nil
}

func (db *datastore) WriteImageFile(w io.Writer, path, hash string) error {
	// Each image has a hash and it's file is stored under a path (one for the
	// images and one for the thumbs), under a folder which begins with the
	// first two letters of the hash.
	f, err := os.Open(filepath.Join(path, hash[0:2], hash))
	if err != nil {
		return fmt.Errorf("could not open image file: %v", err)
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
			return
		}
	}()

	r := bufio.NewReader(f)
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

func (db *datastore) GetRatedImages(username string) ([]shimmie.RatedImage, error) {
	rows, err := db.Query(imageGetSafeBustedQuery, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		img    shimmie.RatedImage
		images []shimmie.RatedImage
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
			&img.Rater,
			&img.RaterIP,
			&img.RateDate,
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

const (

	// imageGetSafeBustedQuery searches score_log for images set as "Safe" ignoring
	// ones from a specific username, extracts the ID of the images and returns
	// them.
	//
	// Warning: MySQL specific query.
	imageGetSafeBustedQuery = `
SELECT
  img.*, rater, rater_ip, rate_date
FROM images as img, (
  SELECT 
    SUBSTRING_INDEX(SUBSTRING_INDEX(message, '#', -1), ' ', 1) AS id,
	score_log.address as rater_ip,
	score_log.username as rater,
	score_log.date_sent as rate_date
  FROM score_log
  WHERE message
  LIKE "%set to: Safe"
  AND username != ?
  ORDER BY date_sent DESC) as safe
WHERE img.id = safe.id
AND rating = 's'
`

	imageGetQuery = `
SELECT * 
FROM images 
WHERE id=?
`
)

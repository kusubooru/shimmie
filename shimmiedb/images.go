package shimmiedb

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kusubooru/shimmie"
)

func (db *DB) CreateImage(ctx context.Context, img shimmie.Image) (int64, error) {
	const query = `
	INSERT INTO images(
		id,
		owner_id,
		owner_ip,
		filename,
		filesize,
		hash,
		ext,
		source,
		width,
		height,
		posted,
		locked
	) values (
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	);`

	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx,
		img.ID,
		img.OwnerID,
		img.OwnerIP,
		img.Filename,
		img.Filesize,
		img.Hash,
		img.Ext,
		img.Source,
		img.Width,
		img.Height,
		img.Posted,
		img.Locked,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *DB) RateImage(id int, rating string) error {
	const query = `
	UPDATE images
	SET rating=?
	WHERE id=?;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := stmt.Close(); err == nil {
			err = cerr
			return
		}
	}()

	res, err := stmt.Exec(rating, id)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()

	return err
}

func (db *DB) GetImage(id int) (*shimmie.Image, error) {
	const query = `
	SELECT *
	FROM images
	WHERE id=?;`

	var (
		img      shimmie.Image
		source   sql.NullString
		parentID sql.NullInt64
		author   sql.NullString
	)
	err := db.QueryRow(query, id).Scan(
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
		&parentID,
		&img.HasChildren,
		&author,
		&img.Notes,
	)
	if err != nil {
		return nil, err
	}
	if source.Valid {
		img.Source = source.String
	}
	if parentID.Valid {
		img.ParentID = parentID.Int64
	}
	if author.Valid {
		img.Author = author.String
	}

	return &img, nil
}

func (db *DB) WriteImageFile(w io.Writer, path, hash string) error {
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
		n, rerr := r.Read(buf)
		if rerr != nil && rerr != io.EOF {
			return rerr
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, werr := w.Write(buf[:n]); werr != nil {
			return werr
		}
	}
	return err
}

func (db *DB) GetRatedImages(username string) ([]shimmie.RatedImage, error) {
	rows, err := db.Query(imageGetRatedQuery, username)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); err == nil {
			err = cerr
			return
		}
	}()

	var images []shimmie.RatedImage
	for rows.Next() {
		var (
			img      shimmie.RatedImage
			source   sql.NullString
			parentID sql.NullInt64
			author   sql.NullString
		)
		err = rows.Scan(
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
			&parentID,
			&img.HasChildren,
			&author,
			&img.Notes,
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
		if parentID.Valid {
			img.ParentID = parentID.Int64
		}
		if author.Valid {
			img.Author = author.String
		}
		images = append(images, img)
	}
	return images, err
}

const (

	// imageGetRatedQuery searches score_log in section "rating" for messages
	// containing "set to: Safe" and extracts the IDs of those images from the
	// log message. It only keeps the latest log ID rows for each extracted
	// image ID. Then it connects the extracted image IDs with images with
	// rating="s" from the "images" table ignoring ones from a specific
	// username (rater) and finally returns those images including the username
	// (rater), the user IP and the date of the original log message.
	//
	// Basically this query allows to find all the images rated as Safe from
	// all users except a specific one while if that specific user rates an
	// image as Safe again (approval), that image won't appear in the results.
	// Since shimmie does not keep a rating history we have to do ugly work
	// using the shimmie log.
	//
	// Warning: MySQL specific query.
	imageGetRatedQuery = `
SELECT
  img.*,
  rater,
  rater_ip,
  rate_date
FROM
  images AS img,
  (
    SELECT
      latest_safe.max_log_id AS max_log_id,
      latest_safe.rated_id AS rated_id,
      log.address AS rater_ip,
      log.username AS rater,
      log.date_sent AS rate_date
    FROM
      (
        SELECT
          MAX(score_log.id) AS max_log_id,
          SUBSTRING_INDEX(SUBSTRING_INDEX(message, '#', - 1), ' ', 1) AS rated_id
        FROM
          score_log
        WHERE
          message LIKE '%set to: Safe'
          AND section = 'rating'
        GROUP BY
          rated_id
        ORDER BY
          max_log_id DESC
      )
      AS latest_safe,
      score_log AS log
    WHERE
      latest_safe.max_log_id = log.id
  )
  AS safe
WHERE
  img.id = safe.rated_id
  AND rating = 's'
  AND rater != ?
`
)

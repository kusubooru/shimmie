package shimmiedb

import (
	"context"
	"time"

	"github.com/kusubooru/shimmie"
)

// CreateTagHistory inserts a new tag history for an image in the db.
func (db *DB) CreateTagHistory(ctx context.Context, th shimmie.TagHistory) (int64, error) {
	if th.DateSet == nil {
		now := time.Now()
		th.DateSet = &now
	}
	const query = `
	INSERT INTO tag_histories(
		id,
		image_id,
		user_id,
		user_ip,
		tags,
		date_set
	) values (
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
		&th.ID,
		&th.ImageID,
		&th.UserID,
		&th.UserIP,
		&th.Tags,
		&th.DateSet,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetImageTagHistory returns the previous tags of an image.
func (db *DB) GetImageTagHistory(imageID int) ([]shimmie.TagHistory, error) {
	rows, err := db.Query(imageTagHistoryGetQuery, imageID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); err == nil {
			err = cerr
			return
		}
	}()

	var ths []shimmie.TagHistory
	for rows.Next() {
		var th shimmie.TagHistory
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

// GetTagHistory returns a tag_history row.
func (db *DB) GetTagHistory(id int) (*shimmie.TagHistory, error) {
	var th shimmie.TagHistory
	err := db.QueryRow(tagHistoryGetQuery, id).Scan(
		&th.ID,
		&th.ImageID,
		&th.UserID,
		&th.UserIP,
		&th.Tags,
		&th.DateSet,
	)
	if err != nil {
		return nil, err
	}
	return &th, err
}

// GetContributedTagHistory returns the latest tag history i.e. tag changes
// that were done by a contributor on an owner's image, per image. It is
// used to fetch data for the "Tag Approval" page.
func (db *DB) GetContributedTagHistory(imageOwnerUsername string) ([]shimmie.ContributedTagHistory, error) {
	rows, err := db.Query(contributedTagHistoryGetQuery, imageOwnerUsername)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); err == nil {
			err = cerr
			return
		}
	}()

	var ths []shimmie.ContributedTagHistory
	for rows.Next() {
		var th shimmie.ContributedTagHistory
		err := rows.Scan(
			&th.ID,
			&th.ImageID,
			&th.OwnerID,
			&th.OwnerName,
			&th.TaggerID,
			&th.TaggerName,
			&th.TaggerIP,
			&th.Tags,
			&th.DateSet,
		)
		if err != nil {
			return nil, err
		}
		ths = append(ths, th)
	}
	return ths, nil
}

const (
	imageTagHistoryGetQuery = `
SELECT tag_histories.*, users.name
FROM tag_histories
JOIN users ON tag_histories.user_id = users.id
WHERE image_id = ?
ORDER BY tag_histories.id DESC
`
	tagHistoryGetQuery = `
SELECT *
FROM tag_histories
WHERE id = ?
`
	// contributedTagHistoryGetQuery performs a "reverse group by" by selecting
	// the max tag_histories ID in a subquery. It does not allow to get the
	// count of tag_histories per image ID as originally planned but it's
	// simpler and faster than queries which include count and group by.
	contributedTagHistoryGetQuery = `
SELECT
    th.id AS id,
    img.id AS image_id,
    owner.id AS owner_id,
    owner.name AS owner_name,
    tagger.id AS tagger_id,
    tagger.name AS tagger_name,
    th.user_ip AS tagger_ip,
    th.tags AS tags,
    th.date_set AS date_set
FROM
    tag_histories th
        JOIN
    images img ON th.image_id = img.id
        JOIN
    users owner ON img.owner_id = owner.id
        JOIN
    users tagger ON th.user_id = tagger.id
WHERE
    th.id = (SELECT
            MAX(id)
        FROM
            tag_histories
        WHERE
            th.image_id = image_id)
        AND th.user_id != img.owner_id
        AND owner.name = ?
ORDER BY date_set DESC
`
)

package store

import "github.com/kusubooru/shimmie"

func (db *datastore) GetImageTagHistory(imageID int) ([]shimmie.TagHistory, error) {
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

func (db *datastore) GetTagHistory(id int) (*shimmie.TagHistory, error) {
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

func (db *datastore) GetContributedTagHistory(imageOwnerUsername string) ([]shimmie.ContributedTagHistory, error) {
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

	var (
		th  shimmie.ContributedTagHistory
		ths []shimmie.ContributedTagHistory
	)
	for rows.Next() {
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

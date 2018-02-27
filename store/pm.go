package store

import (
	"fmt"

	"github.com/kusubooru/shimmie"
)

func (db *datastore) CreatePM(pm *shimmie.PM) error {
	if pm == nil {
		return fmt.Errorf("cannot create nil private message")
	}
	var query = `
	INSERT INTO private_message (
      from_id,
      from_ip,
      to_id,
      sent_date,
      subject,
      message,
      is_read
	)
	VALUES ( ?, ?, ?, ?, ?, ?, ? )
	`
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

	isRead := "N"
	if pm.IsRead {
		isRead = "Y"
	}
	res, err := stmt.Exec(
		pm.FromID,
		pm.FromIP,
		pm.ToID,
		pm.SentDate,
		pm.Subject,
		pm.Message,
		isRead,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	pm.ID = id

	return nil
}

func (db *datastore) GetPMs(from, to string, choice shimmie.PMChoice) ([]*shimmie.PM, error) {
	var query = `
SELECT 
    from_user.name from_user, to_user.name to_user, pm.*
FROM
    private_message pm
        JOIN
    users from_user ON pm.from_id = from_user.id
        JOIN
    users to_user ON pm.to_id = to_user.id
WHERE 1=1
	`

	// gather filters
	var m = make(map[string]interface{})
	if from != "" {
		m["from_user.name"] = from
	}
	if to != "" {
		m["to_user.name"] = to
	}
	switch choice {
	case shimmie.PMRead:
		m["pm.is_read"] = "Y"
	case shimmie.PMUnread:
		m["pm.is_read"] = "N"
	default:
	}

	var values []interface{}
	var filters string
	for k, v := range m {
		filters += fmt.Sprintf(" AND %s = ? ", k)
		values = append(values, v)
	}

	// make query
	rows, err := db.Query(query+filters, values...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); err == nil {
			err = cerr
			return
		}
	}()

	var pms []*shimmie.PM
	var isRead string
	for rows.Next() {
		pm := new(shimmie.PM)
		err = rows.Scan(
			&pm.FromUser,
			&pm.ToUser,
			&pm.ID,
			&pm.FromID,
			&pm.FromIP,
			&pm.ToID,
			&pm.SentDate,
			&pm.Subject,
			&pm.Message,
			&isRead,
		)
		if err != nil {
			return nil, err
		}
		if isRead == "Y" {
			pm.IsRead = true
		}
		pms = append(pms, pm)
	}
	return pms, err
}

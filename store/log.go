package store

import (
	"fmt"
	"time"

	"github.com/kusubooru/shimmie"
)

func (db *datastore) LogRating(imgID int, imgRating, username, userIP string) error {
	rating := shimmie.ImageRating(imgRating)
	msg := fmt.Sprintf("Rating for Image #%d set to: %v", imgID, rating)

	_, err := db.Log("rating", username, userIP, 20, msg)

	return err
}

func (db *datastore) Log(section, username, address string, priority int, message string) (*shimmie.SCoreLog, error) {
	stmt, err := db.Prepare(scoreLogInsertStmt)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	res, err := stmt.Exec(now.Format(time.RFC3339), section, username, address, priority, message)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	log := &shimmie.SCoreLog{
		ID:       id,
		DateSent: &now,
		Section:  section,
		Username: username,
		Address:  address,
		Priority: priority,
		Message:  message,
	}
	return log, nil
}

const scoreLogInsertStmt = `
INSERT score_log
SET
  date_sent=?,
  section=?,
  username=?,
  address=?,
  priority=?,
  message=?
`

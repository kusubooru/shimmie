package shimmie

import (
	"errors"
	"strings"
	"time"
)

const (
	imageRatingSafe         = "Safe"
	imageRatingQuestionable = "Questionable"
	imageRatingExplicit     = "Explicit"
)

// Errors returned by Verify.
var (
	ErrWrongCredentials = errors.New("wrong username or password")
	ErrNotFound         = errors.New("entry not found")
)

// ImageRating converts rating letters to full words.
//
//		s -> Safe
//		q -> Questionable
//		e -> Explicit
//
// If another value except (s, q, e) is given, then it returns that value as it
// is.
func ImageRating(rating string) string {
	switch rating {
	case "s":
		return imageRatingSafe
	case "q":
		return imageRatingQuestionable
	case "e":
		return imageRatingExplicit
	}
	return rating
}

// UserGetter represents a type that can get users from the db.
type UserGetter interface {
	GetUserByName(username string) (*User, error)
}

// Shimmie represents an installed shimmie2 project.
type Shimmie struct {
	ImagePath string
	ThumbPath string
	User      UserGetter
}

// SCoreLog represents a log message in the shimmie log that is stored in the
// table "score_log".
type SCoreLog struct {
	ID       int64
	DateSent *time.Time
	Section  string
	Username string
	Address  string
	Priority int
	Message  string
}

// RatedImage represents a shimmie image that also carries information about
// who rated it and when.
type RatedImage struct {
	Image
	Rater    string
	RaterIP  string
	RateDate *time.Time
}

// RateDateFormat returns the RateDate as UTC with Mon 02 Jan 2006 15:04:05 MST
// format.
func (ri RatedImage) RateDateFormat() string {
	return ri.RateDate.UTC().Format("Mon 02 Jan 2006 15:04:05 MST")
}

// Image represents a shimmie image.
type Image struct {
	ID           int64
	OwnerID      int64
	OwnerIP      string
	Filename     string
	Filesize     int
	Hash         string
	Ext          string
	Source       string
	Width        int
	Height       int
	Posted       *time.Time
	Locked       string
	NumericScore int
	Rating       string
	Favorites    int
	ParentID     int64
	HasChildren  bool
	Author       string
	Notes        int
}

// User represents a shimmie user.
type User struct {
	ID       int64
	Name     string
	Pass     string
	JoinDate *time.Time
	Admin    string
	Email    string
	Class    string
}

// Common holds common configuration values.
type Common struct {
	Title       string
	AnalyticsID string
	Description string
	Keywords    string
}

// SiteTitle returns the Title capitalized.
func (c Common) SiteTitle() string {
	return strings.Title(c.Title)
}

// TagHistory holds previous tags for an image.
type TagHistory struct {
	ID      int64
	ImageID int64
	UserID  int64
	UserIP  string
	Tags    string
	DateSet *time.Time
	// Name of the user who did the edit.
	Name string
}

// ContributedTagHistory holds previous tags for an image that were set by
// contributors.
type ContributedTagHistory struct {
	ID         int
	ImageID    int
	OwnerID    int
	OwnerName  string
	TaggerID   int
	TaggerName string
	TaggerIP   string
	Tags       string
	DateSet    *time.Time
}

// Alias is an alias of an old tag to a new tag.
type Alias struct {
	OldTag string
	NewTag string
}

// Tag is an image's tag.
type Tag struct {
	ID    int
	Tag   string
	Count int
}

// Autocomplete is the result of searching into tags and tag alias to give
// autocomplete suggestions.
type Autocomplete struct {
	Old   string `json:"old"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// PMChoice allows to choose between read and unread private messages.
type PMChoice int

// Possible private message choices.
const (
	PMAny PMChoice = iota
	PMRead
	PMUnread
)

// PM is a private message exchanged between users.
type PM struct {
	FromUser string    `json:"from_user"`
	ToUser   string    `json:"to_user"`
	ID       int64     `json:"id"`
	FromID   int64     `json:"from_id"`
	FromIP   string    `json:"from_ip"`
	ToID     int64     `json:"to_id"`
	SentDate time.Time `json:"sent_date"`
	Subject  string    `json:"subject"`
	Message  string    `json:"message"`
	IsRead   bool      `json:"is_read"`
}

// UserScore can be used to hold user scores like who has uploaded the most
// images and who has edited the most tags.
type UserScore struct {
	Score    int        `json:"score"`
	ID       int64      `json:"id"`
	Name     string     `json:"name"`
	JoinDate *time.Time `json:"join_date"`
	Email    string     `json:"email"`
	Class    string     `json:"class"`
}

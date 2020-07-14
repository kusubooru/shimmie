package shimmie

import (
	"database/sql"
	"errors"
	"io"
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

// Shimmie represents an installed shimmie2 project.
type Shimmie struct {
	ImagePath string
	ThumbPath string
	Store
}

// New creates a new Shimmie by providing a store with the database
// configuration and the paths of the images and thumbs.
func New(imgPath, thumbPath string, s Store) *Shimmie {
	return &Shimmie{ImagePath: imgPath, ThumbPath: thumbPath, Store: s}
}

// Store describes all the operations that need to access database storage.
type Store interface {
	// SQLDB returns the encapsulated *sql.DB, mostly used for testing.
	SQLDB() *sql.DB

	// GetUser gets a user by ID.
	GetUser(userID int64) (*User, error)
	// GetUserByName gets a user by unique username.
	GetUserByName(username string) (*User, error)
	// CreateUser creates a new user and returns their ID.
	CreateUser(*User) error
	// DeleteUser deletes a user based on their ID.
	DeleteUser(int64) error
	// CountUsers returns how many user entries exist in the database.
	CountUsers() (int, error)
	// GetAllUsers returns user entries of the database based on a limit and
	// an offset. If limit < 0, CountUsers will also be executed to get the
	// maximum limit and return all user entries. Offset still works in this
	// case. For example, assuming 10 entries, GetAllUsers(-1, 0), will return
	// all 10 entries and GetAllUsers(-1, 8) will return the last 2 entries.
	GetAllUsers(limit, offset int) ([]User, error)

	// GetConfig gets shimmie config values.
	GetConfig(keys ...string) (map[string]string, error)

	// GetCommon gets common configuration values.
	GetCommon() (*Common, error)

	// GetRatedImages returns all the images that have been rated as safe
	// ignoring the ones from username.
	GetRatedImages(username string) ([]RatedImage, error)
	// GetImage gets a shimmie Image metadata (not it's bytes).
	GetImage(id int) (*Image, error)
	// RateImage sets the rating for an image.
	RateImage(id int, rating string) error
	// WriteImageFile reads a shimmie image file (image or thumb) which exists
	// under a path and has a hash and then writes to w.
	WriteImageFile(w io.Writer, path, hash string) error

	// Log stores a message on score_log table.
	Log(section, username, address string, priority int, message string) (*SCoreLog, error)
	// LogRating logs when an image rating is set.
	LogRating(imgID int, rating, username, userIP string) error

	// GetImageTagHistory returns the previous tags of an image.
	GetImageTagHistory(imageID int) ([]TagHistory, error)
	// GetTagHistory returns a tag_history row.
	GetTagHistory(imageID int) (*TagHistory, error)
	// GetContributedTagHistory returns the latest tag history i.e. tag changes
	// that were done by a contributor on an owner's image, per image. It is
	// used to fetch data for the "Tag Approval" page.
	GetContributedTagHistory(imageOwnerUsername string) ([]ContributedTagHistory, error)

	// GetAlias returns an alias based on its old tag.
	GetAlias(oldTag string) (*Alias, error)
	// CreateAlias creates a new alias.
	CreateAlias(alias *Alias) error
	// DeleteAlias deletes an alias based on its old tag.
	DeleteAlias(oldTag string) error
	// CountAlias returns how many alias entries exist in the database.
	CountAlias() (int, error)
	// GetAllAlias returns alias entries of the database based on a limit and
	// an offset. If limit < 0, CountAlias will also be executed to get the
	// maximum limit and return all alias entries. Offset still works in this
	// case. For example, assuming 10 entries, GetAllAlias(-1, 0), will return
	// all 10 entries and GetAllAlias(-1, 8) will return the last 2 entries.
	GetAllAlias(limit, offset int) ([]Alias, error)
	// FindAlias returns all alias matching an oldTag or a newTag or both.
	FindAlias(oldTag, newTag string) ([]Alias, error)

	// Verify compares the provided username and password with the username and
	// password hash stored in the shimmie database.
	Verify(username, password string) (*User, error)

	CreateTag(*Tag) error
	DeleteTag(name string) error
	GetTag(name string) (*Tag, error)
	GetAllTags(limit, offset int) ([]*Tag, error)

	// Autocomplete searches tags and tag alias for a term and returns
	// suggestions tags to be used for a UI autocomplete.
	Autocomplete(q string, limit, offset int) ([]*Autocomplete, error)

	// CreatePM inserts a new private message.
	CreatePM(pm *PM) error
	// GetPMs returns the private messages exchanged from a user to
	// another user. The arguments from and to are user names and either or
	// both can be left empty.
	GetPMs(from, to string, choice PMChoice) ([]*PM, error)

	// Close closes the connection with the database.
	Close() error
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
	ID           int
	OwnerID      int
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
	ID      int
	ImageID int
	UserID  int
	UserIP  string
	Tags    string
	DateSet *time.Time
	Name    string
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

type PMChoice int

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

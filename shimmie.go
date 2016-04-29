package shimmie

import (
	"io"
	"strings"
	"time"
)

type Shimmie struct {
	ImagePath string
	ThumbPath string
	Store
}

func New(imgPath, thumbPath string, s Store) *Shimmie {
	return &Shimmie{ImagePath: imgPath, ThumbPath: thumbPath, Store: s}
}

// Store describes all the operations that need to access database storage.
type Store interface {
	// GetUser gets a user by unique username.
	GetUser(username string) (*User, error)

	// GetConfig gets shimmie config values.
	GetConfig(keys ...string) (map[string]string, error)

	// GetCommon gets common configuration values.
	GetCommon() (*Common, error)

	// GetSafeBustedImages returns all the images that have been rated as safe
	// ignoring the ones from username.
	GetSafeBustedImages(username string) ([]RatedImage, error)
	// GetImage gets a shimmie Image metadata (not it's bytes).
	GetImage(id int) (*Image, error)
	// WriteImageFile reads a shimmie image file (image or thumb) which exists
	// under a path and has a hash and then writes to w.
	WriteImageFile(w io.Writer, path, hash string) error
}

// RatedImage represents a shimmie image that also carries information about
// who rated it and when.
type RatedImage struct {
	Image
	Rater    string
	RaterIP  string
	RateDate *time.Time
}

// User represents a shimmie image.
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
}

// User represents a shimmie user.
type User struct {
	ID       int
	Name     string
	Pass     string
	JoinDate *time.Time
	Admin    string
	Email    string
	Class    string
}

// CommonConf holds common configuration values.
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

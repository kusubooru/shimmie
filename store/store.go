package store

import (
	"errors"

	"github.com/kusubooru/shimmie"

	"golang.org/x/net/context"
)

const key = "shimmiestore"

// NewContext returns a new Context carrying store.
func NewContext(ctx context.Context, store Store) context.Context {
	return context.WithValue(ctx, key, store)
}

// FromContext extracts the store from ctx, if present.
func FromContext(ctx context.Context) (Store, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the Store type assertion returns ok=false for nil.
	s, ok := ctx.Value(key).(Store)
	return s, ok
}

// Store describes all the operations that need to access database storage.
type Store interface {
	// GetUser gets a user by unique username.
	GetUser(username string) (*shimmie.User, error)

	// GetConfig gets shimmie config values.
	GetConfig(keys ...string) (map[string]string, error)

	// GetSafeBustedImages returns all the images that have been rated as safe
	// ignoring the ones from username.
	GetSafeBustedImages(username string) ([]shimmie.Image, error)
}

// GetSafeBustedImages returns all the images that have been rated as safe
// ignoring the ones from username.
func GetSafeBustedImages(ctx context.Context, username string) ([]shimmie.Image, error) {
	s, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("no store in context")
	}
	return s.GetSafeBustedImages(username)
}

// GetConfig gets shimmie config values.
func GetConfig(ctx context.Context, keys ...string) (map[string]string, error) {
	s, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("no store in context")
	}
	return s.GetConfig(keys...)
}

// GetUser gets a user by unique username.
func GetUser(ctx context.Context, username string) (*shimmie.User, error) {
	s, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("no store in context")
	}
	return s.GetUser(username)
}

package store_test

import (
	"database/sql"
	"log"
	"reflect"
	"testing"

	"github.com/kusubooru/shimmie"
	"github.com/kusubooru/shimmie/store"
)

func TestUser(t *testing.T) {
	schema := setup()
	defer teardown(schema)
	shim := store.Open(*driverName, *dataSourceName)
	defer func() {
		if cerr := shim.Close(); cerr != nil {
			log.Println("failed to close connection")
		}
	}()

	username := "john"
	password := "1234"
	u := &shimmie.User{
		Name:  username,
		Pass:  password,
		Email: "john@doe.com",
	}

	// Create a user.
	err := shim.CreateUser(u)
	if err != nil {
		t.Fatalf("CreateUser(%q) returned err: %v", u, err)
	}
	expectedID := int64(1)
	if got, want := u.ID, expectedID; got != want {
		t.Errorf("CreateUser(%q) -> user.Id = %d, want %d", u, got, want)
	}
	if got, want := u.Pass, shimmie.Hash(password); got != want {
		t.Errorf("CreateUser(%q) -> user.Pass = %q, want %q", u, got, want)
	}
	if got, want := u.Class, "user"; got != want {
		t.Errorf("CreateUser(%q) -> user.Class = %q, want %q", u, got, want)
	}

	// Attempt to get created user and compare.
	got, err := shim.GetUser(username)
	if err != nil {
		t.Fatalf("GetUser(%q) returned err: %v", username, err)
	}
	if want := u; !reflect.DeepEqual(got, want) {
		t.Errorf("GetUser(%q) -> user =\n%#v, want\n%#v", username, got, want)
	}

	// Delete created user.
	if err := shim.DeleteUser(expectedID); err != nil {
		t.Errorf("DeleteUser(%d) returned err: %v", expectedID, err)
	}

	// Attempt to get user again and expect no rows err.
	_, err = shim.GetUser(username)
	if got, want := err, sql.ErrNoRows; got != want {
		t.Errorf("GetUser(%q) after delete returned err = %v, want %v", username, got, want)
	}
}

package store_test

import (
	"database/sql"
	"fmt"
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
	got, err := shim.GetUserByName(username)
	if err != nil {
		t.Fatalf("GetUserByName(%q) returned err: %v", username, err)
	}
	if want := u; !reflect.DeepEqual(got, want) {
		t.Errorf("GetUserByName(%q) -> user =\n%#v, want\n%#v", username, got, want)
	}

	// Also get user by ID just to test the method.
	got, err = shim.GetUser(expectedID)
	if err != nil {
		t.Fatalf("GetUser(%d) returned err: %v", expectedID, err)
	}
	if want := u; !reflect.DeepEqual(got, want) {
		t.Errorf("GetUser(%d) -> user =\n%#v, want\n%#v", expectedID, got, want)
	}

	// Delete created user.
	if err := shim.DeleteUser(expectedID); err != nil {
		t.Errorf("DeleteUser(%d) returned err: %v", expectedID, err)
	}

	// Attempt to get user again and expect no rows err.
	_, err = shim.GetUserByName(username)
	if got, want := err, sql.ErrNoRows; got != want {
		t.Errorf("GetUserByName(%q) after delete returned err = %v, want %v", username, got, want)
	}
}

func TestGetAllUsers(t *testing.T) {
	schema := setup()
	defer teardown(schema)
	shim := store.Open(*driverName, *dataSourceName)
	defer func() {
		if cerr := shim.Close(); cerr != nil {
			log.Println("failed to close connection")
		}
	}()

	pass, max := "123", 10
	for i := 0; i < max; i++ {
		u := &shimmie.User{
			Name: fmt.Sprintf("user%d", i),
			Pass: pass,
		}
		err := shim.CreateUser(u)
		if err != nil {
			t.Fatalf("CreateUser(%q) returned err: %v", u, err)
		}
	}

	var getAllUserTests = []struct {
		limit   int
		offset  int
		wantLen int
	}{
		// Get all users with limit and offset.
		{limit: 5, offset: 0, wantLen: 5},
		// Get all users in the database by providing a negative limit.
		{limit: -1, offset: 8, wantLen: 2},
		// Get all users with offset that exceeds the number of entries.
		{limit: 10, offset: 20, wantLen: 0},
	}

	for _, tt := range getAllUserTests {
		limit, offset := tt.limit, tt.offset
		users, err := shim.GetAllUsers(limit, offset)
		if err != nil {
			t.Fatalf("GetAllUsers(%d, %d) returned err: %v", limit, offset, err)
		}
		if got, want := len(users), tt.wantLen; got != want {
			t.Errorf("GetAllUsers(%d, %d) -> len(users) = %d, want %d", limit, offset, got, want)
		}
	}
}

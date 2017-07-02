package store_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/kusubooru/shimmie"
	"github.com/kusubooru/shimmie/store"
)

func TestVerify(t *testing.T) {
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

	// Verify success case.
	got, err := shim.Verify(username, password)
	if err != nil {
		t.Fatalf("Verify(%q, %q) returned err: %v", username, password, err)
	}
	if want := u; !reflect.DeepEqual(got, want) {
		t.Errorf("Verify(%q, %q) -> user =\n%#v, want\n%#v", username, password, got, want)
	}

	// Verify wrong password case.
	password = "wrongpassword"
	_, err = shim.Verify(username, password)
	if err == nil {
		t.Fatalf("Verify(%q, %q) expected to return err", username, password)
	}
	if got, want := err, shimmie.ErrWrongCredentials; got != want {
		t.Errorf("Verify(%q, %q) -> err =\n%#v, want\n%#v", username, password, got, want)
	}

	// Verify user not found case.
	username = "nonexistentuser"
	password = "wrongpassword"
	_, err = shim.Verify(username, password)
	if err == nil {
		t.Fatalf("Verify(%q, %q) expected to return err", username, password)
	}
	if got, want := err, shimmie.ErrNotFound; got != want {
		t.Errorf("Verify(%q, %q) -> err =\n%#v, want\n%#v", username, password, got, want)
	}
}

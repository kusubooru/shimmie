package store_test

import (
	"reflect"
	"testing"

	"github.com/kusubooru/shimmie"
)

func TestVerify(t *testing.T) {
	shim, schema := setup(t)
	defer teardown(t, shim, schema)

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

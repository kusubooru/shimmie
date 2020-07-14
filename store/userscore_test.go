package store_test

import (
	"testing"

	"github.com/kusubooru/shimmie"
)

func TestMostImageUploads(t *testing.T) {
	shim, schema := setup(t)
	defer teardown(t, shim, schema)

	var users = []shimmie.User{
		{Name: "bob", Pass: "bob123"},
		{Name: "ann", Pass: "ann123"},
		{Name: "zoe", Pass: "zoe123"},
	}
	for _, u := range users {
		if err := shim.CreateUser(&u); err != nil {
			t.Fatalf("CreateUser(%q) returned err: %v", u, err)
		}
	}

	// TODO(jin): Create images

	// TODO(jin): Change setup to return concrete type.

	// TODO(jin): Run tests in docker-compose
	//
	// https://www.ardanlabs.com/blog/2019/03/integration-testing-in-go-executing-tests-with-docker.html

	// TODO(jin): skip db tests when go test -short

	// TODO(jin): run all tests using docker-compose and go run make.go -test

	t.Errorf("must implement create images")
	t.Errorf("must add images to users")
}

package shimmiedb_test

import (
	"context"
	"testing"

	"github.com/kusubooru/shimmie"
)

func TestCreateImage(t *testing.T) {
	shim, schema := setup(t)
	defer teardown(t, shim, schema)

	u := shimmie.User{Name: "bob", Pass: "bob123"}
	if err := shim.CreateUser(&u); err != nil {
		t.Fatalf("CreateUser(%#v) returned err: %v", u, err)
	}

	ctx := context.Background()
	img := shimmie.Image{OwnerID: u.ID}
	t.Logf("CreateImage(img) should succeed: img=%#v", img)
	id, err := shim.CreateImage(ctx, img)
	if err != nil {
		t.Errorf("CreateImage(img) returned error: %v", err)
	}
	if id == 0 {
		t.Error("CreateImage(img) should return a non-zero id")
	}
}

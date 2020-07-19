package store_test

import (
	"context"
	"testing"

	"github.com/kusubooru/shimmie"
)

func TestCreateTagHistory(t *testing.T) {
	shim, schema := setup(t)
	defer teardown(t, shim, schema)

	u := shimmie.User{Name: "bob"}
	if err := shim.CreateUser(&u); err != nil {
		t.Fatalf("CreateUser(%#v) returned err: %v", u, err)
	}

	ctx := context.Background()
	img := shimmie.Image{OwnerID: u.ID}
	imgID, err := shim.CreateImage(ctx, img)
	if err != nil {
		t.Fatalf("CreateImage() for user %q returned err: %v", u.Name, err)
	}

	th := shimmie.TagHistory{ImageID: imgID, UserID: u.ID}
	t.Logf("CreateTagHistory(th) should succeed: th=%#v", th)
	id, err := shim.CreateTagHistory(ctx, th)
	if err != nil {
		t.Errorf("CreateTagHistory(th) returned error: %v", err)
	}
	if id == 0 {
		t.Error("CreateTagHistory(th) should return a non-zero id")
	}
}

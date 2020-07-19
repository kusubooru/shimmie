package shimmiedb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kusubooru/shimmie"
)

func TestMostImageUploads(t *testing.T) {
	shim, schema := setup(t)
	defer teardown(t, shim, schema)

	fixtures := []struct {
		user   shimmie.User
		images []shimmie.Image
	}{
		{
			user: shimmie.User{Name: "bob"},
			images: []shimmie.Image{
				{OwnerID: 1, Hash: "0"},
			},
		},
		{
			user: shimmie.User{Name: "ann"},
			images: []shimmie.Image{
				{OwnerID: 2, Hash: "1"},
				{OwnerID: 2, Hash: "2"},
			},
		},
		{
			user: shimmie.User{Name: "zoe"},
			images: []shimmie.Image{
				{OwnerID: 3, Hash: "3"},
				{OwnerID: 3, Hash: "4"},
				{OwnerID: 3, Hash: "5"},
			},
		},
	}
	ctx := context.Background()
	t.Log("After inserting:")
	for _, f := range fixtures {
		if err := shim.CreateUser(&f.user); err != nil {
			t.Fatalf("CreateUser(%v) returned err: %v", f.user, err)
		}
		t.Logf("User %q with images:", f.user.Name)
		for _, img := range f.images {
			id, err := shim.CreateImage(ctx, img)
			if err != nil {
				t.Errorf("CreateImage(%#v) returned error: %v", img, err)
			}
			t.Logf("|-> Image id=%d", id)
		}
	}

	score, err := shim.MostImageUploads(10)
	if err != nil {
		t.Fatalf("MostImageUploads() returned err: %v", err)
	}
	want := []shimmie.UserScore{
		{Score: 3, Name: "zoe"},
		{Score: 2, Name: "ann"},
		{Score: 1, Name: "bob"},
	}
	t.Log("and calling MostImageUploads() the user scores should be:")
	for i, s := range want {
		t.Logf("score[%d]-> Score: %d, Name: %q", i, s.Score, s.Name)
	}
	t.Logf("we got:")
	for i, s := range score {
		prefix := fmt.Sprintf("score[%d]->", i)
		t.Logf("%s Score: %d, Name: %q", prefix, s.Score, s.Name)
		testUserScore(t, score[i], want[i], prefix)
	}
}

func testUserScore(t *testing.T, got, want shimmie.UserScore, prefix string) {
	t.Helper()
	if got, want := got.Score, want.Score; got != want {
		t.Errorf("%s Score: %d, want: %d", prefix, got, want)
	}
	if got, want := got.Name, want.Name; got != want {
		t.Errorf("%s Name: %q, want: %q", prefix, got, want)
	}
}

func TestMostTagEdits(t *testing.T) {
	shim, schema := setup(t)
	defer teardown(t, shim, schema)

	fixtures := []struct {
		user         shimmie.User
		images       []shimmie.Image
		tagHistories []shimmie.TagHistory
	}{
		{
			user: shimmie.User{Name: "bob"},
			images: []shimmie.Image{
				{OwnerID: 1, Hash: "0"},
			},
			tagHistories: []shimmie.TagHistory{
				{UserID: 1, ImageID: 1},
			},
		},
		{
			user: shimmie.User{Name: "ann"},
			images: []shimmie.Image{
				{OwnerID: 2, Hash: "1"},
			},
			tagHistories: []shimmie.TagHistory{
				{UserID: 2, ImageID: 2},
				{UserID: 2, ImageID: 2},
			},
		},
		{
			user: shimmie.User{Name: "zoe"},
			images: []shimmie.Image{
				{OwnerID: 3, Hash: "2"},
			},
			tagHistories: []shimmie.TagHistory{
				{UserID: 3, ImageID: 3},
				{UserID: 3, ImageID: 3},
				{UserID: 3, ImageID: 3},
			},
		},
	}
	ctx := context.Background()
	t.Log("After inserting:")
	for _, f := range fixtures {
		if err := shim.CreateUser(&f.user); err != nil {
			t.Fatalf("CreateUser(%v) returned err: %v", f.user, err)
		}
		t.Logf("User %q with images and tag histories:", f.user.Name)
		for _, img := range f.images {
			id, err := shim.CreateImage(ctx, img)
			if err != nil {
				t.Errorf("CreateImage(%#v) returned error: %v", img, err)
			}
			t.Logf("|-> Image id=%d", id)
		}
		for _, th := range f.tagHistories {
			id, err := shim.CreateTagHistory(ctx, th)
			if err != nil {
				t.Errorf("CreateTagHistory(%#v) returned error: %v", th, err)
			}
			t.Logf("|-> TagHistory id=%d", id)
		}
	}

	score, err := shim.MostTagEdits(10)
	if err != nil {
		t.Fatalf("MostTagEdits() returned err: %v", err)
	}
	want := []shimmie.UserScore{
		{Score: 3, Name: "zoe"},
		{Score: 2, Name: "ann"},
		{Score: 1, Name: "bob"},
	}
	t.Log("and calling MostTagEdits() the user scores should be:")
	for i, s := range want {
		t.Logf("score[%d]-> Score: %d, Name: %q", i, s.Score, s.Name)
	}
	t.Logf("we got:")
	for i, s := range score {
		prefix := fmt.Sprintf("score[%d]->", i)
		t.Logf("%s Score: %d, Name: %q", prefix, s.Score, s.Name)
		testUserScore(t, score[i], want[i], prefix)
	}
}

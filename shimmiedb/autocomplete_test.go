package shimmiedb_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/kusubooru/shimmie"
)

func TestAutocomplete(t *testing.T) {
	shim, schema := setup(t)
	defer teardown(t, shim, schema)

	// Test that searching with empty query returns empty results.
	empty := ""
	emptyAutocomplete, err := shim.Autocomplete(empty, 10, 0)
	if err != nil {
		t.Fatalf("Autocomplete(%q) returned err: %v", empty, err)
	}
	if got, want := len(emptyAutocomplete), 0; got != want {
		t.Fatalf("Autocomplete(%q) returned %d results but expected %d instead", empty, got, want)
	}

	tagName := "character:chun-li"
	chunLiTag := &shimmie.Tag{Tag: tagName, Count: 5}
	// Create a tag.
	if err := shim.CreateTag(chunLiTag); err != nil {
		t.Fatalf("CreateTag(%q) returned err: %v", chunLiTag, err)
	}

	alias := &shimmie.Alias{
		OldTag: "chun-li",
		NewTag: tagName,
	}
	// Create an Alias for the previous tag.
	if err := shim.CreateAlias(alias); err != nil {
		t.Fatalf("CreateAlias(%q) returned err: %v", alias, err)
	}

	chunTag := &shimmie.Tag{Tag: "chun", Count: 1}
	// Create another unrelated tag.
	if err := shim.CreateTag(chunTag); err != nil {
		t.Fatalf("CreateTag(%q) returned err: %v", chunTag, err)
	}

	// Do an autocomplete query for a term that should include both the tag
	// with its alias and the unrelated tag.
	q := "chun"
	tags, err := shim.Autocomplete(q, 10, 0)
	if err != nil {
		t.Fatalf("Autocomplete(%q) returned err: %v", q, err)
	}

	expected := []*shimmie.Autocomplete{
		{Old: "chun-li", Name: "character:chun-li", Count: 5},
		{Old: "", Name: "chun", Count: 1},
	}
	if got, want := len(tags), len(expected); got != want {
		t.Errorf("Autocomplete(%q) returned %d results but expected to return %d instead", q, got, want)
	}
	if got, want := tags, expected; !reflect.DeepEqual(got, want) {
		t.Errorf("Autocomplete(%q) -> %#+v, want %#+v", q, got, want)
		data, _ := json.Marshal(got)
		fmt.Println(string(data))
		data, _ = json.Marshal(want)
		fmt.Println(string(data))
	}
}

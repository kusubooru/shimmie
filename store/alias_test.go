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

func TestAlias(t *testing.T) {
	schema := setup()
	defer teardown(schema)
	shim := store.Open(*driverName, *dataSourceName)
	defer func() {
		if cerr := shim.Close(); cerr != nil {
			log.Println("failed to close connection")
		}
	}()

	oldTag := "old_tag"
	alias := &shimmie.Alias{
		OldTag: oldTag,
		NewTag: "new_tag",
	}
	// Create an Alias.
	err := shim.CreateAlias(alias)
	if err != nil {
		t.Fatalf("CreateAlias(%q) returned err: %v", alias, err)
	}

	anotherAlias := &shimmie.Alias{
		OldTag: oldTag,
		NewTag: "some_other_new_tag",
	}
	// Attempt to create a new alias that has the same old tag and expect
	// error.
	err = shim.CreateAlias(anotherAlias)
	if err == nil {
		t.Fatalf("CreateAlias(%q) must return err because it has\n"+
			"  the same old tag as %q which already exists in the database.", anotherAlias, alias)
	}

	// Get successfully created alias.
	got, err := shim.GetAlias(oldTag)
	if err != nil {
		t.Fatalf("GetAlias(%q) returned err: %v", oldTag, err)
	}
	if want := alias; !reflect.DeepEqual(got, want) {
		t.Errorf("GetAlias(%q) -> %q, want %q", oldTag, got, want)
	}

	// Count created alias and find only 1.
	count, err := shim.CountAlias()
	if err != nil {
		t.Errorf("CountAlias() returned err: %v", err)
	}
	if got, want := count, 1; got != want {
		t.Errorf("CountAlias() -> count = %d, want %d", got, want)
	}

	// Deleted created alias.
	if err := shim.DeleteAlias(oldTag); err != nil {
		t.Errorf("DeleteAlias(%d) returned err: %v", oldTag, err)
	}

	// Attempt to get alias again and expect no rows err.
	_, err = shim.GetAlias(oldTag)
	if got, want := err, sql.ErrNoRows; got != want {
		t.Errorf("GetAlias(%q) after delete returned err = %v, want %v", oldTag, got, want)
	}

}

func TestGetAllAlias(t *testing.T) {
	schema := setup()
	defer teardown(schema)
	shim := store.Open(*driverName, *dataSourceName)
	defer func() {
		if cerr := shim.Close(); cerr != nil {
			log.Println("failed to close connection")
		}
	}()

	newTag, max := "old_tag", 10
	for i := 0; i < max; i++ {
		a := &shimmie.Alias{
			OldTag: fmt.Sprintf("old_tag%d", i),
			NewTag: newTag,
		}
		err := shim.CreateAlias(a)
		if err != nil {
			t.Fatalf("CreateAlias(%q) returned err: %v", a, err)
		}
	}

	var getAllAliasTests = []struct {
		limit   int
		offset  int
		wantLen int
	}{
		// Get all alias with limit and offset.
		{limit: 5, offset: 0, wantLen: 5},
		// Get all alias in the database by providing a negative limit.
		{limit: -1, offset: 8, wantLen: 2},
		// Get all alias with offset that exceeds the number of entries.
		{limit: 10, offset: 20, wantLen: 0},
	}

	for _, tt := range getAllAliasTests {
		limit, offset := tt.limit, tt.offset
		alias, err := shim.GetAllAlias(limit, offset)
		if err != nil {
			t.Fatalf("GetAllAlias(%d, %d) returned err: %v", limit, offset, err)
		}
		if got, want := len(alias), tt.wantLen; got != want {
			t.Errorf("GetAllAlias(%d, %d) -> len(alias) = %d, want %d", limit, offset, got, want)
		}
	}
}

func TestFindAlias(t *testing.T) {
	schema := setup()
	defer teardown(schema)
	shim := store.Open(*driverName, *dataSourceName)
	defer func() {
		if cerr := shim.Close(); cerr != nil {
			log.Println("failed to close connection")
		}
	}()

	alias := []shimmie.Alias{
		{NewTag: "character:sarah_fortune", OldTag: "character:miss_fortune"},
		{NewTag: "character:sarah_fortune", OldTag: "miss_fortune"},
		{NewTag: "character:sarah_fortune", OldTag: "miss_fortune_the_bounty_hunter"},
		{NewTag: "character:sarah_fortune", OldTag: "Miss_Forturne_(lol)"},
		{NewTag: "character:sarah_fortune", OldTag: "sarah_fortune"},
	}
	for _, a := range alias {
		err := shim.CreateAlias(&a)
		if err != nil {
			t.Fatalf("CreateAlias(%q) returned err: %v", a, err)
		}
	}

	tests := []struct {
		oldTag  string
		newTag  string
		matches int
	}{
		{"", "fortune", 5},
		{"miss", "", 4}, // 4 because it also finds Miss.
		{"character", "character", 1},
		{"", "", 5},
	}

	for _, tt := range tests {
		res, err := shim.FindAlias(tt.oldTag, tt.newTag)
		if err != nil {
			t.Fatalf("FindAlias(%q, %q) returned err: %v", tt.oldTag, tt.newTag, err)
		}
		if got, want := len(res), tt.matches; got != want {
			t.Errorf("FindAlias(%q, %q) len(result) = %d, want %d", tt.oldTag, tt.newTag, got, want)
		}
	}
}

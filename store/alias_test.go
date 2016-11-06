package store_test

import (
	"database/sql"
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

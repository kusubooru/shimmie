package store_test

import (
	"context"
	"flag"
	"testing"

	"github.com/kusubooru/shimmie/store"
)

const (
	defaultDataSource = "shimmie:shimmie@(127.0.0.1:3306)/shimmie?parseTime=true&multiStatements=true"
)

var (
	testDataSource = flag.String("datasource", defaultDataSource, "database data source used for tests")
)

func setup(t *testing.T) (*store.Datastore, *store.Schema) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}

	schema, err := store.NewSchemer(*testDataSource, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := schema.DB.Ping(); err != nil {
		t.Fatalf("make sure database is up: use docker-compose up -d: %v", err)
	}
	err = schema.Create()
	if err != nil {
		t.Fatalf("failed to create schema using datasource %s: %v", *testDataSource, err)
	}

	shim, err := store.Open(*testDataSource, 10)
	if err != nil {
		t.Fatalf("failed to connect using datasource %s: %v", *testDataSource, err)
	}
	return shim, schema
}

func teardown(t *testing.T, shim *store.Datastore, schema *store.Schema) {
	if err := schema.TruncateTables(context.Background()); err != nil {
		t.Errorf("error truncating tables: %v", err)
	}
}

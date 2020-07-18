package store_test

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/kusubooru/shimmie/store"
)

const (
	defaultUsername     = "shimmie"
	defaultPassword     = "shimmie"
	defaultHost         = "127.0.0.1"
	defaultPort         = "3306"
	defaultTestDatabase = "shimmie"
	defaultDriver       = "mysql"
)

var (
	username          = flag.String("u", defaultUsername, "connect with this username to run the tests")
	password          = flag.String("p", defaultPassword, "connect with this password to run the tests")
	host              = flag.String("host", defaultHost, "connect on this host to run the tests")
	port              = flag.String("port", defaultPort, "connect on this port to run the tests")
	testDBName        = flag.String("dbname", defaultTestDatabase, "create and drop this database on each test")
	driverName        = flag.String("driver", defaultDriver, "database driver")
	defaultDataSource = fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", *username, *password, *host, *port, *testDBName)
	dataSourceName    = flag.String("datasource", defaultDataSource, "database data source")
)

func setup(t *testing.T) (*store.Datastore, *store.Schema) {
	schema := store.NewSchemer(*driverName, *username, *password, *host, *port, *testDBName, 0)
	err := schema.DB.Ping()
	if err != nil {
		t.Fatalf("make sure database is up: use docker-compose up -d: %v", err)
	}

	err = schema.Create(*testDBName)
	if err != nil {
		t.Fatalf("failed to create schema for %s: %v", *testDBName, err)
	}
	shim := store.Open(*dataSourceName, 10)
	return shim, schema
}

func teardown(t *testing.T, shim *store.Datastore, schema *store.Schema) {
	if err := schema.TruncateTables(context.Background(), *testDBName); err != nil {
		t.Errorf("error truncating tables: %v", err)
	}
}

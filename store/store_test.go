package store_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/kusubooru/shimmie"
	"github.com/kusubooru/shimmie/store"
)

const (
	defaultUsername     = "kusubooru"
	defaultPassword     = "kusubooru"
	defaultHost         = "localhost"
	defaultPort         = "3306"
	defaultTestDatabase = "kusubooru_test"
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

func setup(t *testing.T) (shimmie.Store, shimmie.Schemer) {
	schema := store.NewSchemer(*driverName, *username, *password, *host, *port)
	err := schema.Create(*testDBName)
	if err != nil {
		t.Fatalf("failed to create schema for %s: %v", *testDBName, err)
	}
	shim := store.Open(*driverName, *dataSourceName)
	return shim, schema
}

func teardown(t *testing.T, shim shimmie.Store, schema shimmie.Schemer) {
	if err := shim.Close(); err != nil {
		t.Errorf("error closing shimmie connection: %v", err)
	}

	if err := schema.Drop(*testDBName); err != nil {
		t.Errorf("error dropping schema: %v", err)
	}
	if err := schema.Close(); err != nil {
		t.Errorf("error closing schema connection: %v", err)
	}
}

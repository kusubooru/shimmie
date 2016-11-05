package store_test

import (
	"flag"
	"fmt"
	"log"

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

func init() {
	flag.Parse()
}

func setup() shimmie.Schemer {
	schema := store.NewSchemer(*driverName, *username, *password, *host, *port)
	err := schema.Create(*testDBName)
	if err != nil {
		log.Printf("failed to create schema for %s: %v", *testDBName, err)
	}
	return schema
}

func teardown(schema shimmie.Schemer) {
	if err := schema.Drop(*testDBName); err != nil {
		log.Println(err)
	}
	if err := schema.Close(); err != nil {
		log.Println(err)
	}
}

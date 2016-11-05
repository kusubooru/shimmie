// +build ignore

package main

import (
	"flag"
	"log"

	"github.com/kusubooru/shimmie/store"
)

const (
	defaultUsername     = "kusubooru"
	defaultPassword     = "kusubooru"
	defaultHost         = "localhost"
	defaultPort         = "3306"
	defaultTestDatabase = "kusubooru_dev"
	defaultDriver       = "mysql"
)

var (
	username   = flag.String("u", defaultUsername, "connect with this username to make the schema")
	password   = flag.String("p", defaultPassword, "connect with this password to make the schema")
	host       = flag.String("host", defaultHost, "connect on this host to make the schema")
	port       = flag.String("port", defaultPort, "connect on this port to make the schema")
	devDBName  = flag.String("dbname", defaultTestDatabase, "make the schema on this database")
	driverName = flag.String("driver", defaultDriver, "database driver")
	dropSchema = flag.Bool("drop", false, "drop schema if true")
)

func main() {
	flag.Parse()

	schema := store.NewSchemer(*driverName, *username, *password, *host, *port)
	defer func() {
		if cerr := schema.Close(); cerr != nil {
			log.Println("Could not close database connection:", cerr)
		}
	}()

	if *dropSchema {
		if err := schema.Drop(*devDBName); err != nil {
			log.Println("Error dropping schema:", err)
			return
		}
	}

	if err := schema.Create(*devDBName); err != nil {
		log.Println("Error creating schema:", err)
		return
	}
}

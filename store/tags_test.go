package store_test

import (
	"fmt"
	"testing"
)

func TestGetContributedHistory(t *testing.T) {
	shim := setup(t)
	defer teardown(t, shim.SQLDB())

	ths, err := shim.GetContributedTagHistory("kusubooru")
	if err != nil {
		t.Error("err is:", err)
	}
	if ths == nil {
		t.Error("ths is nil")
	}
	if len(ths) == 0 {
		t.Error("ths are empty")
	}
	fmt.Println(ths)
}

// Setup the test environment.
//func setup() (*DB, error) {
//	//err := withTestDB()
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	// testOptions is a global in this case, but you could easily
//	// create one per-test
//	db, err := openDB(*dbDriver, *dbConfig)
//	if err != nil {
//		return nil, err
//	}
//
//	// Loads our test schema
//	db.MustLoad()
//	return db, nil
//}

// Create our test database.
//func withTestDB() error {
//	db, err := Open()
//	if err != nil {
//		return err
//	}
//	defer db.Close()
//
//	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", testOptions.Name))
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

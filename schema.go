package shimmie

// Schemer describes operations that can be done on the database schema. Those
// operations include creating and dropping the schema for a specific database
// name.
type Schemer interface {
	Create(dbName string) error
	Drop(dbName string) error
	Close() error
}

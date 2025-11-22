package database

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
)

// BadgerDB wraps Badger database with basic operations
type BadgerDB struct {
	db *badger.DB
}

// NewBadgerDB creates a new Badger database instance
func NewBadgerDB(path string) (*BadgerDB, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil // Disable badger logging to reduce noise

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger database: %w", err)
	}

	bdb := &BadgerDB{db: db}

	log.Printf("📁 Badger database opened: %s", path)
	return bdb, nil
}

// Close closes the database
func (bdb *BadgerDB) Close() error {
	return bdb.db.Close()
}

// DB returns the underlying badger database instance
func (bdb *BadgerDB) DB() *badger.DB {
	return bdb.db
}


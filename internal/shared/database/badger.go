package database

import (
	"encoding/json"
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

// Generic key-value operations
func (bdb *BadgerDB) Set(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (bdb *BadgerDB) Get(key string, dest any) error {
	return bdb.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, dest)
		})
	})
}

func (bdb *BadgerDB) Delete(key string) error {
	return bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// BatchSet sets multiple key-value pairs in a single transaction
func (bdb *BadgerDB) BatchSet(items map[string]any) error {
	return bdb.db.Update(func(txn *badger.Txn) error {
		for key, value := range items {
			data, err := json.Marshal(value)
			if err != nil {
				return err
			}

			if err := txn.Set([]byte(key), data); err != nil {
				return err
			}
		}
		return nil
	})
}

// Keys returns all keys with a given prefix
func (bdb *BadgerDB) Keys(prefix string) ([]string, error) {
	var keys []string

	err := bdb.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
			item := it.Item()
			key := string(item.Key())
			keys = append(keys, key)
		}
		return nil
	})

	return keys, err
}

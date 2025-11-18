package db

import (
	"github.com/dgraph-io/badger/v4"
)

type DB struct {
	*badger.DB
}

func Init() (*DB, error) {
	opts := badger.DefaultOptions("./data")
	opts.Logger = nil // Disable badger logs
	
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func (db *DB) Set(key, value []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (db *DB) Get(key []byte) ([]byte, error) {
	var value []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		
		return item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
	})
	
	return value, err
}
package time

import (
	"encoding/json"
	"fmt"
	"time"

	"datastar-go/internal/shared/types"

	"github.com/dgraph-io/badger/v4"
)

type Repository struct {
	db *badger.DB
}

func NewRepository(db *badger.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(entry *types.TimeEntry) error {
	key := fmt.Sprintf("time_entry:%s", entry.ID)
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (r *Repository) Get(id string) (*types.TimeEntry, error) {
	key := fmt.Sprintf("time_entry:%s", id)
	var entry types.TimeEntry

	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &entry)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (r *Repository) GetByUserID(userID string) ([]*types.TimeEntry, error) {
	prefix := "time_entry:"
	var entries []*types.TimeEntry

	err := r.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
			item := it.Item()

			err := item.Value(func(val []byte) error {
				var entry types.TimeEntry
				if err := json.Unmarshal(val, &entry); err != nil {
					return err
				}

				if entry.UserID == userID {
					entries = append(entries, &entry)
				}
				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil
	})

	return entries, err
}

func (r *Repository) GetActiveTimer(userID string) (*types.TimeEntry, error) {
	entries, err := r.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsRunning {
			return entry, nil
		}
	}

	return nil, nil
}

func (r *Repository) Delete(id string) error {
	key := fmt.Sprintf("time_entry:%s", id)
	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (r *Repository) GetByProjectIDAndDateRange(projectID string, startDate, endDate time.Time) ([]*types.TimeEntry, error) {
	var entries []*types.TimeEntry
	err := r.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("time_entry:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var entry types.TimeEntry
				if err := json.Unmarshal(val, &entry); err != nil {
					return err
				}
				if entry.ProjectID == projectID &&
					entry.StartTime.After(startDate) &&
					entry.StartTime.Before(endDate) {
					entries = append(entries, &entry)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return entries, err
}

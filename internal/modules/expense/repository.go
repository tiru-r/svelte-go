package expense

import (
	"encoding/json"

	"datastar-go/internal/shared/types"

	"github.com/dgraph-io/badger/v4"
)

type Repository struct {
	db *badger.DB
}

func NewRepository(db *badger.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(expense *types.Expense) error {
	return r.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(expense)
		if err != nil {
			return err
		}
		return txn.Set([]byte("expense:"+expense.ID), data)
	})
}

func (r *Repository) GetByID(id string) (*types.Expense, error) {
	var expense types.Expense
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("expense:" + id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &expense)
		})
	})
	return &expense, err
}

func (r *Repository) GetByUserID(userID string) ([]*types.Expense, error) {
	var expenses []*types.Expense
	err := r.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("expense:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var expense types.Expense
				if err := json.Unmarshal(val, &expense); err != nil {
					return err
				}
				if expense.UserID == userID {
					expenses = append(expenses, &expense)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return expenses, err
}

func (r *Repository) GetByProjectID(projectID string) ([]*types.Expense, error) {
	var expenses []*types.Expense
	err := r.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("expense:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var expense types.Expense
				if err := json.Unmarshal(val, &expense); err != nil {
					return err
				}
				if expense.ProjectID == projectID {
					expenses = append(expenses, &expense)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return expenses, err
}

func (r *Repository) Update(expense *types.Expense) error {
	return r.Create(expense)
}

func (r *Repository) Delete(id string) error {
	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("expense:" + id))
	})
}

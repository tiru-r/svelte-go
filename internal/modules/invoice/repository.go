package invoice

import (
	"encoding/json"

	"svelte-go/internal/shared/types"

	"github.com/dgraph-io/badger/v4"
)

type Repository struct {
	db *badger.DB
}

func NewRepository(db *badger.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(invoice *types.Invoice) error {
	return r.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(invoice)
		if err != nil {
			return err
		}
		return txn.Set([]byte("invoice:"+invoice.ID), data)
	})
}

func (r *Repository) GetByID(id string) (*types.Invoice, error) {
	var invoice types.Invoice
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("invoice:" + id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &invoice)
		})
	})
	return &invoice, err
}

func (r *Repository) GetByUserID(userID string) ([]*types.Invoice, error) {
	var invoices []*types.Invoice
	err := r.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("invoice:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var invoice types.Invoice
				if err := json.Unmarshal(val, &invoice); err != nil {
					return err
				}
				if invoice.UserID == userID {
					invoices = append(invoices, &invoice)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return invoices, err
}

func (r *Repository) GetByClientID(clientID string) ([]*types.Invoice, error) {
	var invoices []*types.Invoice
	err := r.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("invoice:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var invoice types.Invoice
				if err := json.Unmarshal(val, &invoice); err != nil {
					return err
				}
				if invoice.ClientID == clientID {
					invoices = append(invoices, &invoice)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return invoices, err
}

func (r *Repository) Update(invoice *types.Invoice) error {
	return r.Create(invoice)
}

func (r *Repository) Delete(id string) error {
	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("invoice:" + id))
	})
}

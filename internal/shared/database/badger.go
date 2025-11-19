package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
	"svelte-go/internal/shared/types"
)

// BadgerDB wraps Badger database with business logic
type BadgerDB struct {
	db          *badger.DB
	TimeEntryRepo *TimeEntryRepository
	ProjectRepo   *ProjectRepository
	ClientRepo    *ClientRepository
	ExpenseRepo   *ExpenseRepository
	InvoiceRepo   *InvoiceRepository
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
	
	// Initialize repositories
	bdb.TimeEntryRepo = &TimeEntryRepository{db: db}
	bdb.ProjectRepo = &ProjectRepository{db: db}  
	bdb.ClientRepo = &ClientRepository{db: db}
	bdb.ExpenseRepo = &ExpenseRepository{db: db}
	bdb.InvoiceRepo = &InvoiceRepository{db: db}
	
	log.Printf("📁 Badger database opened: %s", path)
	return bdb, nil
}

// Close closes the database
func (bdb *BadgerDB) Close() error {
	return bdb.db.Close()
}

// Repository patterns for different entities

// TimeEntryRepository handles time entry storage
type TimeEntryRepository struct {
	db *badger.DB
}

// NewTimeEntryRepository creates a new repository
func (bdb *BadgerDB) TimeEntries() *TimeEntryRepository {
	return &TimeEntryRepository{db: bdb.db}
}

// Save stores a time entry
func (r *TimeEntryRepository) Save(entry *types.TimeEntry) error {
	key := fmt.Sprintf("time_entry:%s", entry.ID)
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

// Get retrieves a time entry by ID
func (r *TimeEntryRepository) Get(id string) (*types.TimeEntry, error) {
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

// GetByUserID gets all time entries for a user
func (r *TimeEntryRepository) GetByUserID(userID string) ([]*types.TimeEntry, error) {
	prefix := fmt.Sprintf("time_entry:")
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

// GetActiveTimer gets active timer for a user
func (r *TimeEntryRepository) GetActiveTimer(userID string) (*types.TimeEntry, error) {
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

// Delete removes a time entry
func (r *TimeEntryRepository) Delete(id string) error {
	key := fmt.Sprintf("time_entry:%s", id)
	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// GetByProjectIDAndDateRange retrieves time entries for a project within date range
func (r *TimeEntryRepository) GetByProjectIDAndDateRange(projectID string, startDate, endDate time.Time) ([]*types.TimeEntry, error) {
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

// ProjectRepository handles project storage
type ProjectRepository struct {
	db *badger.DB
}

func (bdb *BadgerDB) Projects() *ProjectRepository {
	return &ProjectRepository{db: bdb.db}
}

func (r *ProjectRepository) Save(project *types.Project) error {
	key := fmt.Sprintf("project:%s", project.ID)
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}

	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (r *ProjectRepository) Get(id string) (*types.Project, error) {
	key := fmt.Sprintf("project:%s", id)
	var project types.Project

	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &project)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return &project, err
}

// Create stores a new project (alias for Save)
func (r *ProjectRepository) Create(project *types.Project) error {
	return r.Save(project)
}

// GetByID retrieves a project by ID (alias for Get)
func (r *ProjectRepository) GetByID(id string) (*types.Project, error) {
	return r.Get(id)
}

// Update updates a project (alias for Save)
func (r *ProjectRepository) Update(project *types.Project) error {
	return r.Save(project)
}

// GetByUserID retrieves projects for a user
func (r *ProjectRepository) GetByUserID(userID string) ([]*types.Project, error) {
	var projects []*types.Project
	err := r.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("project:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var project types.Project
				if err := json.Unmarshal(val, &project); err != nil {
					return err
				}
				if project.UserID == userID {
					projects = append(projects, &project)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return projects, err
}

// GetByClientID retrieves projects for a client
func (r *ProjectRepository) GetByClientID(clientID string) ([]*types.Project, error) {
	var projects []*types.Project
	err := r.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("project:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var project types.Project
				if err := json.Unmarshal(val, &project); err != nil {
					return err
				}
				if project.ClientID == clientID {
					projects = append(projects, &project)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return projects, err
}

// ClientRepository handles client storage  
type ClientRepository struct {
	db *badger.DB
}

func (bdb *BadgerDB) Clients() *ClientRepository {
	return &ClientRepository{db: bdb.db}
}

func (r *ClientRepository) Save(client *types.Client) error {
	key := fmt.Sprintf("client:%s", client.ID)
	data, err := json.Marshal(client)
	if err != nil {
		return err
	}

	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (r *ClientRepository) Get(id string) (*types.Client, error) {
	key := fmt.Sprintf("client:%s", id)
	var client types.Client

	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &client)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return &client, err
}

// Create stores a new client (alias for Save)
func (r *ClientRepository) Create(client *types.Client) error {
	return r.Save(client)
}

// GetByID retrieves a client by ID (alias for Get)
func (r *ClientRepository) GetByID(id string) (*types.Client, error) {
	return r.Get(id)
}

// Update updates a client (alias for Save)
func (r *ClientRepository) Update(client *types.Client) error {
	return r.Save(client)
}

// GetByUserID retrieves clients for a user
func (r *ClientRepository) GetByUserID(userID string) ([]*types.Client, error) {
	var clients []*types.Client
	err := r.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("client:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var client types.Client
				if err := json.Unmarshal(val, &client); err != nil {
					return err
				}
				if client.UserID == userID {
					clients = append(clients, &client)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return clients, err
}

// UserRepository handles user storage
type UserRepository struct {
	db *badger.DB
}

func (bdb *BadgerDB) Users() *UserRepository {
	return &UserRepository{db: bdb.db}
}

func (r *UserRepository) Save(user *types.User) error {
	key := fmt.Sprintf("user:%s", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (r *UserRepository) Get(id string) (*types.User, error) {
	key := fmt.Sprintf("user:%s", id)
	var user types.User

	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return &user, err
}

// Generic key-value operations
func (bdb *BadgerDB) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (bdb *BadgerDB) Get(key string, dest interface{}) error {
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
func (bdb *BadgerDB) BatchSet(items map[string]interface{}) error {
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

// ExpenseRepository handles expense storage
type ExpenseRepository struct {
	db *badger.DB
}

// Create stores a new expense
func (r *ExpenseRepository) Create(expense *types.Expense) error {
	return r.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(expense)
		if err != nil {
			return err
		}
		return txn.Set([]byte("expense:"+expense.ID), data)
	})
}

// GetByID retrieves an expense by ID
func (r *ExpenseRepository) GetByID(id string) (*types.Expense, error) {
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

// GetByUserID retrieves expenses for a user
func (r *ExpenseRepository) GetByUserID(userID string) ([]*types.Expense, error) {
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

// GetByProjectID retrieves expenses for a project
func (r *ExpenseRepository) GetByProjectID(projectID string) ([]*types.Expense, error) {
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

// Update updates an expense
func (r *ExpenseRepository) Update(expense *types.Expense) error {
	return r.Create(expense) // Same as create for key-value store
}

// Delete deletes an expense
func (r *ExpenseRepository) Delete(id string) error {
	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("expense:" + id))
	})
}

// InvoiceRepository handles invoice storage
type InvoiceRepository struct {
	db *badger.DB
}

// Create stores a new invoice
func (r *InvoiceRepository) Create(invoice *types.Invoice) error {
	return r.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(invoice)
		if err != nil {
			return err
		}
		return txn.Set([]byte("invoice:"+invoice.ID), data)
	})
}

// GetByID retrieves an invoice by ID
func (r *InvoiceRepository) GetByID(id string) (*types.Invoice, error) {
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

// GetByUserID retrieves invoices for a user
func (r *InvoiceRepository) GetByUserID(userID string) ([]*types.Invoice, error) {
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

// GetByClientID retrieves invoices for a client
func (r *InvoiceRepository) GetByClientID(clientID string) ([]*types.Invoice, error) {
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

// Update updates an invoice
func (r *InvoiceRepository) Update(invoice *types.Invoice) error {
	return r.Create(invoice) // Same as create for key-value store
}

// Delete deletes an invoice
func (r *InvoiceRepository) Delete(id string) error {
	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("invoice:" + id))
	})
}
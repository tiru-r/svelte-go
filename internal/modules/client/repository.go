package client

import (
	"encoding/json"
	"fmt"

	"svelte-go/internal/shared/types"

	"github.com/dgraph-io/badger/v4"
)

type ClientRepository struct {
	db *badger.DB
}

type ProjectRepository struct {
	db *badger.DB
}

func NewClientRepository(db *badger.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func NewProjectRepository(db *badger.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
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

func (r *ClientRepository) Create(client *types.Client) error {
	return r.Save(client)
}

func (r *ClientRepository) GetByID(id string) (*types.Client, error) {
	return r.Get(id)
}

func (r *ClientRepository) Update(client *types.Client) error {
	return r.Save(client)
}

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

func (r *ProjectRepository) Create(project *types.Project) error {
	return r.Save(project)
}

func (r *ProjectRepository) GetByID(id string) (*types.Project, error) {
	return r.Get(id)
}

func (r *ProjectRepository) Update(project *types.Project) error {
	return r.Save(project)
}

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

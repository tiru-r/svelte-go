package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
)

type Repository struct {
	db *badger.DB
}

func NewRepository(db *badger.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true

	return r.db.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("user:%s", user.ID)
		emailKey := fmt.Sprintf("user_email:%s", user.Email)
		usernameKey := fmt.Sprintf("user_username:%s", user.Username)

		_, err := txn.Get([]byte(emailKey))
		if err == nil {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}

		_, err = txn.Get([]byte(usernameKey))
		if err == nil {
			return fmt.Errorf("user with username %s already exists", user.Username)
		}

		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		if err := txn.Set([]byte(key), data); err != nil {
			return err
		}

		if err := txn.Set([]byte(emailKey), []byte(user.ID)); err != nil {
			return err
		}

		return txn.Set([]byte(usernameKey), []byte(user.ID))
	})
}

func (r *Repository) GetUserByID(id string) (*User, error) {
	var user User
	err := r.db.View(func(txn *badger.Txn) error {
		key := fmt.Sprintf("user:%s", id)
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})
	})
	return &user, err
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	var userID string
	err := r.db.View(func(txn *badger.Txn) error {
		emailKey := fmt.Sprintf("user_email:%s", email)
		item, err := txn.Get([]byte(emailKey))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			userID = string(val)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return r.GetUserByID(userID)
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var userID string
	err := r.db.View(func(txn *badger.Txn) error {
		usernameKey := fmt.Sprintf("user_username:%s", username)
		item, err := txn.Get([]byte(usernameKey))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			userID = string(val)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return r.GetUserByID(userID)
}

func (r *Repository) UpdateUser(user *User) error {
	user.UpdatedAt = time.Now()

	return r.db.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("user:%s", user.ID)
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return txn.Set([]byte(key), data)
	})
}

func (r *Repository) CreateSession(session *Session) error {
	session.ID = uuid.New().String()
	session.CreatedAt = time.Now()
	session.IsActive = true

	return r.db.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("session:%s", session.ID)
		userSessionKey := fmt.Sprintf("user_session:%s", session.UserID)
		
		data, err := json.Marshal(session)
		if err != nil {
			return err
		}

		if err := txn.Set([]byte(key), data); err != nil {
			return err
		}

		return txn.Set([]byte(userSessionKey), []byte(session.ID))
	})
}

func (r *Repository) GetSessionByToken(token string) (*Session, error) {
	var session Session
	err := r.db.View(func(txn *badger.Txn) error {
		key := fmt.Sprintf("session_token:%s", token)
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		var sessionID string
		err = item.Value(func(val []byte) error {
			sessionID = string(val)
			return nil
		})
		if err != nil {
			return err
		}

		sessionKey := fmt.Sprintf("session:%s", sessionID)
		sessionItem, err := txn.Get([]byte(sessionKey))
		if err != nil {
			return err
		}

		return sessionItem.Value(func(val []byte) error {
			return json.Unmarshal(val, &session)
		})
	})
	return &session, err
}

func (r *Repository) InvalidateSession(sessionID string) error {
	return r.db.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("session:%s", sessionID)
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		var session Session
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &session)
		})
		if err != nil {
			return err
		}

		session.IsActive = false
		data, err := json.Marshal(session)
		if err != nil {
			return err
		}

		return txn.Set([]byte(key), data)
	})
}

func (r *Repository) InvalidateUserSessions(userID string) error {
	return r.db.Update(func(txn *badger.Txn) error {
		userSessionKey := fmt.Sprintf("user_session:%s", userID)
		item, err := txn.Get([]byte(userSessionKey))
		if err != nil {
			return err
		}

		var sessionID string
		err = item.Value(func(val []byte) error {
			sessionID = string(val)
			return nil
		})
		if err != nil {
			return err
		}

		return r.InvalidateSession(sessionID)
	})
}
package repo

import (
	"context"
	"encoding/json"
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/ningzio/geminal/internal"
)

var _ internal.Repository = (*Repository)(nil)

func NewRepository() (*Repository, error) {
	db, err := badger.Open(badger.DefaultOptions(".geminal.db"))
	if err != nil {
		return nil, err
	}
	return &Repository{
		db: db,
	}, nil
}

type Repository struct {
	db *badger.DB
}

// DeleteConversation implements internal.Repository.
func (repo *Repository) DeleteConversation(ctx context.Context, chatID string) error {
	err := repo.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(chatStoreKey(chatID))
	})
	if err != nil {
		return err
	}
	return nil
}

var chatStoreKeyPrefix = []byte("conversation:")

func chatStoreKey(chatID string) []byte {
	return append(chatStoreKeyPrefix, []byte(chatID)...)
}

// GetConversationByChatID implements internal.Repository.
func (repo *Repository) GetConversationByChatID(ctx context.Context, chatID string) (*internal.Conversation, error) {
	var conversation internal.Conversation
	err := repo.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(chatStoreKey(chatID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &conversation)
		})
	})
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

// LoadHistory implements internal.Repository.
func (repo *Repository) LoadHistory(ctx context.Context) ([]*internal.Conversation, error) {
	var conversations []*internal.Conversation
	err := repo.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		opts.Prefix = chatStoreKeyPrefix
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var conversation internal.Conversation
			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &conversation)
			})
			if err != nil {
				return err
			}
			conversations = append(conversations, &conversation)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

// SaveConversation implements internal.Repository.
func (repo *Repository) SaveConversation(ctx context.Context, conversation *internal.Conversation) error {
	err := repo.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(conversation)
		if err != nil {
			return err
		}
		return txn.Set(chatStoreKey(conversation.ChatID), data)
	})

	if err != nil {
		return fmt.Errorf("saving conversation: %w", err)
	}
	return nil
}

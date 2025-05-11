package store

import "errors"

// ErrKeyNotFound is returned when a key is not found in the store
var ErrKeyNotFound = errors.New("key not found")

// KeyStore defines the interface for a key-value store
type KeyStore interface {
	// Get retrieves a value by its key
	Get(key string) (interface{}, error)

	// Set stores a value with the given key
	Set(key string, value interface{}) error

	// Delete removes a value by its key
	Delete(key string) error

	// Clear removes all entries from the store
	Clear() error

	// Keys returns all keys in the store
	Keys() []string
}

// InMemoryKeyStore implements KeyStore interface using in-memory map
type InMemoryKeyStore struct {
	store map[string]interface{}
}

// NewInMemoryKeyStore creates a new InMemoryKeyStore instance
func NewInMemoryKeyStore() *InMemoryKeyStore {
	return &InMemoryKeyStore{
		store: make(map[string]interface{}),
	}
}

func (s *InMemoryKeyStore) Get(key string) (interface{}, error) {
	value, exists := s.store[key]
	if !exists {
		return nil, ErrKeyNotFound
	}
	return value, nil
}

func (s *InMemoryKeyStore) Set(key string, value interface{}) error {
	s.store[key] = value
	return nil
}

func (s *InMemoryKeyStore) Delete(key string) error {
	delete(s.store, key)
	return nil
}

func (s *InMemoryKeyStore) Clear() error {
	s.store = make(map[string]interface{})
	return nil
}

func (s *InMemoryKeyStore) Keys() []string {
	keys := make([]string, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k)
	}
	return keys
}

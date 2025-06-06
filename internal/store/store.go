package store

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

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

// DiskKeyStore implements KeyStore interface with persistent storage
type DiskKeyStore struct {
	filepath string
	data     map[string]interface{}
	mu       sync.RWMutex
}

// NewDiskKeyStore creates a new DiskKeyStore instance
func NewDiskKeyStore(filepath string) (*DiskKeyStore, error) {
	store := &DiskKeyStore{
		filepath: filepath,
		data:     make(map[string]interface{}),
	}

	if err := store.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return store, nil
}

func (s *DiskKeyStore) load() error {
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.data)
}

func (s *DiskKeyStore) save() error {
	data, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filepath, data, 0644)
}

func (s *DiskKeyStore) Get(key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if value, ok := s.data[key]; ok {
		return value, nil
	}
	return nil, ErrKeyNotFound
}

func (s *DiskKeyStore) Set(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return s.save()
}

func (s *DiskKeyStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	return s.save()
}

func (s *DiskKeyStore) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]interface{})
	return s.save()
}

func (s *DiskKeyStore) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}
